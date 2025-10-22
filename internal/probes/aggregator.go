// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/result"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/report/status"

	"github.com/gofrs/uuid/v5"
)

type Aggregator struct {
	// Identification of the karaoke
	Repository string      `json:"repository"`
	Songname   string      `json:"songname"`
	Kid        uuid.UUID   `json:"kid"`
	CreatedAt  time.Time   `json:"created-at"`
	ModifiedAt time.Time   `json:"modified-at"`
	Year       json.Number `json:"year"`
	// `Probes` report direct features of the karaoke based on metadata, lyrics, etc.
	// They can be used to detect common mistakes.
	Probes  []probe.Probe            `json:"-"`
	Reports map[string]report.Report `json:"reports"`
	Stats   *Stats                   `json:"statistics"`
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		Probes:  EnabledProbes(),
		Reports: make(map[string]report.Report),
		Stats:   &Stats{},
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
	a.Repository = karaJson.Data.Repository
	a.Songname = karaJson.Data.Songname
	a.Kid = karaJson.Data.Kid
	a.CreatedAt = karaJson.Data.CreatedAt
	a.ModifiedAt = karaJson.Data.ModifiedAt
	a.Year = karaJson.Data.Year
	a.Stats.Reset()
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
				switch r.r.Status() {
				case status.Completed:
					switch r.r.Result() {
					case result.Passed:
						a.Stats.Passed += 1
					case result.Failed:
						switch r.r.Severity() {
						case severity.Info:
							a.Stats.FailedInfo += 1
						case severity.Warning:
							a.Stats.FailedWarning += 1
						case severity.Critical:
							a.Stats.FailedCritical += 1
						}
					}
				case status.Aborted:
					a.Stats.Aborted += 1
				case status.Skipped:
					a.Stats.Skipped += 1
				}

			}
		}
		return nil
	}
}
