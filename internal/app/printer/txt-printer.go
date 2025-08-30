// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package printer

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/louisroyer/km-probe/internal/app/ansi"
	"github.com/louisroyer/km-probe/internal/probes"
	"github.com/louisroyer/km-probe/internal/probes/report/result"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/report/status"

	"github.com/moby/term"
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
	size, err := term.GetWinsize(os.Stdout.Fd())
	if err != nil {
		return err
	}

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
		fmt.Printf("%s", ansi.Link(u, a.Songname+" ["+a.Kid.String()+"] ("+a.Repository+")")) // TODO: split on 2 lines if too long for terminal
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
		fmt.Printf("  %s: ", k) // TODO: alignment
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
		wrap := int(max(1, min(size.Width-4, 120))) // wrap at 120 col. max
		if msg := r.Message(); msg != "" {
			cursor := 0
			for next := min(wrap, len(msg)); cursor < len(msg); next = min(next+wrap, len(msg)) {
				// avoid splitting inside a word
				for pos := next; pos < len(msg); pos++ {
					// make eventual firsts space characters of the next line part of this line
					// (they will be trimmed)
					if string(msg[pos]) != " " {
						next = pos + 1
						break
					}
				}
				if next < len(msg) {
					for pos := next - 1; pos > cursor; pos-- {
						// find the last space character of the line,
						// and make it the end of what we will print (it will be trimmed)
						if string(msg[pos]) == " " {
							next = pos
							break
						}
					}
				}
				fmt.Printf("    %s\n", strings.TrimSpace(msg[cursor:next]))
				cursor = next
			}
		}
	}
	fmt.Printf("\n") // empty line to separate aggregators
	return nil
}
