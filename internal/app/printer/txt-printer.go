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

	"github.com/karaoke-tools/km-probe/internal/app/ansi"
	"github.com/karaoke-tools/km-probe/internal/probes"
	"github.com/karaoke-tools/km-probe/internal/probes/report/result"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/report/status"

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
	builder := &strings.Builder{}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ready:
		defer p.setReady()
		defer p.aggregatorPool.Put(a)
		builder.Reset()
		if err := p.encodeAggregator(ctx, a, builder); err != nil {
			return err
		}
	}
	return nil

}

func (p *TxtPrinter) encodeAggregator(ctx context.Context, a *probes.Aggregator, builder *strings.Builder) error {
	size, err := term.GetWinsize(os.Stdout.Fd())
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s [%s] (%s)", a.Songname, a.Kid.String(), a.Repository)
	cursor := 0
	wrap := int(max(1, size.Width))
	for next := min(wrap, len(msg)); cursor < len(msg); next = min(next+wrap, len(msg)) {
		if p.Color {
			if a.Stats.FailedCritical > 0 {
				builder.WriteString(ansi.Red)
			} else if a.Stats.FailedWarning > 0 {
				builder.WriteString(ansi.Yellow)
			} else if a.Stats.Passed > 0 {
				builder.WriteString(ansi.Green)
			} else {
				builder.WriteString(ansi.Blue)
			}
		}

		if p.Hyperlink {
			builder.WriteString(ansi.LinkStart)
			u, err := url.JoinPath(p.BaseUri, a.Kid.String())
			if err != nil {
				return err
			}
			builder.WriteString(u)
			builder.WriteString(ansi.LinkMiddle)
		}

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
		builder.WriteString(strings.TrimSpace(msg[cursor:next]))
		cursor = next

		if p.Hyperlink {
			builder.WriteString(ansi.LinkEnd)
		}
		if p.Color {
			builder.WriteString(ansi.Reset)
		}

		fmt.Println(builder.String())
		builder.Reset()
	}

	first := true
	if a.Stats.Passed > 0 {
		builder.WriteString(fmt.Sprintf("Passed: %d", a.Stats.Passed))
		first = false
	}
	if a.Stats.FailedCritical > 0 {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("Failed (critical): %d", a.Stats.FailedCritical))
		first = false
	}
	if a.Stats.FailedWarning > 0 {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("Failed (warning): %d", a.Stats.FailedWarning))
		first = false
	}
	if a.Stats.FailedInfo > 0 {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("Info: %d", a.Stats.FailedInfo))
		first = false
	}
	if a.Stats.Skipped > 0 {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("Skipped: %d", a.Stats.Skipped))
		first = false
	}
	if a.Stats.Aborted > 0 {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("Aborted: %d", a.Stats.Aborted))
		first = false
	}
	msg = builder.String()
	builder.Reset()
	wrap = int(max(1, size.Width-4)) // 2 is size of start of line "  > "
	cursor = 0
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
		fmt.Printf("  > %s\n", strings.TrimSpace(msg[cursor:next]))
		cursor = next
	}

	for k, r := range a.Reports {
		if r.Status() != status.Completed || r.Result() != result.Failed {
			continue
		}
		builder.WriteString("  ")
		if p.Color {
			switch r.Severity() {
			case severity.Critical:
				builder.WriteString(ansi.Red)
			case severity.Warning:
				builder.WriteString(ansi.Yellow)
			case severity.Info:
				builder.WriteString(ansi.Blue)
			}
		}
		builder.WriteString(k)
		builder.WriteString(" ")
		if r.Result() == result.Failed {
			if r.Severity() != severity.Info {
				builder.WriteString(fmt.Sprintf("[%s (%s)]", r.Result(), r.Severity()))
			} else {
				builder.WriteString("[info]")
			}
		}

		if p.Color {
			builder.WriteString(ansi.Reset)
		}
		fmt.Println(builder.String())
		builder.Reset()

		wrap := int(max(1, size.Width-4)) // 4 is number of spaces at start of line in final formating
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
