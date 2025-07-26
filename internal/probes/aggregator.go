// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"context"
	"fmt"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type Aggregator struct {
	Name    string
	Reports map[string]report.Report
	Probes  []probe.Probe
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
		return nil
	}
}

func (a *Aggregator) String() string {
	ret := []string{fmt.Sprintf("name: %s", a.Name)}
	for k, v := range a.Reports {
		ret = append(ret, fmt.Sprintf("- %s: %s", k, v))
	}
	ret = append(ret, fmt.Sprintf("- probably-good-first-contribution: %t", a.SuitableFirstContribution()))
	return strings.Join(ret, "\n")
}

func (a *Aggregator) SuitableFirstContribution() bool {
	issue_cnt := 0
	if !a.Reports["style-single-white"].Result() || !a.Reports["style-black-border"].Result() {
		issue_cnt += 1
	}
	if !a.Reports["eol-punctuation"].Result() {
		issue_cnt += 1
	}
	if !a.Reports["double-consonnant"].Result() {
		issue_cnt += 1
	}

	if a.Reports["automation"].Result() && a.Reports["live-download"].Result() && a.Reports["resolution"].Result() && issue_cnt == 1 {
		return true
	}
	return false
}
