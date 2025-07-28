// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type Aggregator struct {
	Name    string                   `json:"name"`
	Reports map[string]report.Report `json:"reports"`
	Probes  []probe.Probe            `json:"-"`
}

func FromKaraJson(ctx context.Context, basedir string, karaJson *karajson.KaraJson, probes *[]probe.NewProbeFunc) (*Aggregator, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		data, err := karadata.FromKaraJson(ctx, basedir, karaJson)
		if err != nil {
			return nil, err
		}
		aggregator := Aggregator{
			Name:    data.KaraJson.Data.Songname,
			Reports: make(map[string]report.Report),
		}
		if probes == nil {
			probes = &defaultProbes
		}
		for _, probe := range *probes {
			aggregator.Probes = append(aggregator.Probes, probe(data))
		}
		return &aggregator, nil
	}
}

type reportWithName struct {
	name string
	r    report.Report
}

func (a *Aggregator) Run(ctx context.Context) error {
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
					if r, err := p.Run(ctx); err == nil {
						ch <- reportWithName{name: p.Name(), r: r}
					} else {
						ch <- reportWithName{name: p.Name(), r: report.Abort()}
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
		// add additional reports
		a.Reports["probably-good-first-contribution"] = a.SuitableFirstContribution()
		return nil
	}
}

func (a *Aggregator) SuitableFirstContribution() report.Report {
	critical := []string{
		"live-download",
		"resolution",
		"automation",
	}
	for _, c := range critical {
		if r, ok := a.Reports[c]; !ok {
			return report.Skip()
		} else if !r.Result() {
			return report.Info(false)
		}
	}

	scoring := [][]string{
		// style issues
		[]string{"style-single-white", "style-black-border"}, // minor issues
		[]string{"resolution"},                               // can imply re-splitting some parts
		// lyrics issues
		[]string{"double-consonnant"},
	}

	badness := 0
	for _, s := range scoring {
		local_badness := 0
		for _, sb := range s {
			if r, ok := a.Reports[sb]; !ok {
				return report.Skip()
			} else if !r.Result() {
				local_badness++
			}
		}
		if local_badness > 0 {
			badness++
		}
	}

	return report.Info(badness == 1)
}
