// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/report"

	"github.com/gofrs/uuid"
)

type Aggregator struct {
	// Identification of the karaoke
	Songname   string      `json:"songname"`
	Kid        uuid.UUID   `json:"kid"`
	CreatedAt  time.Time   `json:"created-at"`
	ModifiedAt time.Time   `json:"modified-at"`
	Year       json.Number `json:"year"`
	// `Probes` report direct features of the karaoke based on metadata, lyrics, etc.
	// They can be used to detect common mistakes.
	Probes  []probe.Probe            `json:"-"`
	Reports map[string]report.Report `json:"reports"`
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		Reports: make(map[string]report.Report),
		Probes:  AvailableProbes(),
	}
}

func (a *Aggregator) Reset(basedir string, karaJson *karajson.KaraJson) {
	// recycle reports & analysis memory
	for _, v := range a.Reports {
		v.Delete()
	}
	// empty the map
	clear(a.Reports)
	if karaJson == nil {
		return
	}
	a.Songname = karaJson.Data.Songname
	a.Kid = karaJson.Data.Kid
	a.CreatedAt = karaJson.Data.CreatedAt
	a.ModifiedAt = karaJson.Data.ModifiedAt
	a.Year = karaJson.Data.Year
}

type reportWithName struct {
	name string
	r    report.Report
}

func (a *Aggregator) Run(ctx context.Context, KaraData *karadata.KaraData) error {
	select {
	// if a.Probes is empty, context would not be checked otherwise
	case <-ctx.Done():
		return ctx.Err()
	default:
		ch := make(chan reportWithName)
		// start probes
		for _, p := range a.Probes {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				go func(ctx context.Context, p probe.Probe, ch chan<- reportWithName) {
					if s, msg, err := p.PreRun(ctx, KaraData); err != nil {
						select {
						case <-ctx.Done():
							return
						case ch <- reportWithName{name: p.Name(), r: report.Abort()}:
							return
						}
					} else if s {
						select {
						case <-ctx.Done():
							return
						case ch <- reportWithName{name: p.Name(), r: report.Skip(msg)}:
							return
						}
					}
					if r, err := p.Run(ctx, KaraData); err == nil {
						select {
						case <-ctx.Done():
							return
						case ch <- reportWithName{name: p.Name(), r: r}:
							return
						}
					} else {
						select {
						case <-ctx.Done():
							return
						case ch <- reportWithName{name: p.Name(), r: report.Abort()}:
							return

						}
					}
				}(ctx, p, ch)
			}
		}
		// get result of probes
		for _, _ = range a.Probes {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case r := <-ch:
				a.Reports[r.name] = r.r
			}
		}
		return nil
	}
}
