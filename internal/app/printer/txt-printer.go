// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package printer

import (
	"context"
	"fmt"
	"net/url"

	"github.com/louisroyer/km-probe/internal/app/ansi"
	"github.com/louisroyer/km-probe/internal/probes"
	"github.com/louisroyer/km-probe/internal/probes/report/result"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/report/status"
)

type TxtPrinter struct {
	*BasePrinter
	Hyperlink bool
	Color     bool
	BaseUri   string
}

func NewTxtPrinter(hyperlink bool, color bool, baseUri string) Printer {
	return &TxtPrinter{
		BasePrinter: NewBasePrinter(),
		Hyperlink:   hyperlink,
		Color:       color,
		BaseUri:     baseUri,
	}
}

func (p *TxtPrinter) Encode(ctx context.Context, a *probes.Aggregator) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ready:
		defer p.setReady()
		defer p.aggregatorPool.Put(a)
		if err := p.encodeAggregator(ctx, a); err != nil {
			return err
		}
	}
	return nil

}

func (p *TxtPrinter) encodeAggregator(ctx context.Context, a *probes.Aggregator) error {

	if p.Hyperlink {
		u, err := url.JoinPath(p.BaseUri, a.Kid.String())
		if err != nil {
			return err
		}
		if p.Color {
			if a.Stats.FailedCritical > 0 {
				fmt.Printf(ansi.Red)
			} else if a.Stats.FailedWarning > 0 {
				fmt.Printf(ansi.Yellow)
			} else if a.Stats.Passed > 0 {
				fmt.Printf(ansi.Green)
			} else {
				fmt.Printf(ansi.Blue)
			}
		}
		fmt.Printf("%s", ansi.Link(u, a.Songname+" ["+a.Kid.String()+"] ("+a.Repository+")"))
		if p.Color {
			fmt.Printf(ansi.Reset)
		}
		fmt.Printf("\n")
	} else {
		fmt.Println(a.Songname)
	}
	// TODO: build this to remove fields with 0
	fmt.Printf("Passed: %d, Failed (critical): %d, Failed (warning): %d, Info: %d, Skipped: %d, Aborted: %d\n", a.Stats.Passed, a.Stats.FailedCritical, a.Stats.FailedWarning, a.Stats.FailedInfo, a.Stats.Skipped, a.Stats.Aborted)
	for k, r := range a.Reports {
		if r.Status() != status.Completed || r.Result() != result.Failed {
			continue
		}
		if p.Color {
			switch r.Severity() {
			case severity.Critical:
				fmt.Printf(ansi.Red)
			case severity.Warning:
				fmt.Printf(ansi.Yellow)
			case severity.Info:
				fmt.Printf(ansi.Blue)
			}
		}
		fmt.Printf("\t%s: ", k) // TODO: alignment
		if r.Result() == result.Failed {
			if r.Severity() != severity.Info {
				fmt.Printf("%s [%s]", r.Result(), r.Severity())
			} else {
				fmt.Printf("[info]")
			}
		}

		if p.Color {
			fmt.Printf(ansi.Reset)
		}
		fmt.Printf("\n")
		if msg := r.Message(); msg != "" {
			fmt.Printf("\t\t%s\n", msg) // TODO: wrap around 120 col. (but don't split inside a word)
		}
	}
	fmt.Printf("\n") // empty line to separate aggregators
	return nil
}
