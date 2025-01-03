// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strconv"
	"strings"
)

type Style struct {
	LyricsFile *Ass
}

func NewStyle(lyrics *Ass) *Style {
	return &Style{LyricsFile: lyrics}
}

func (p *Style) Run(ctx context.Context) (*Report, error) {
	report := NewReport()
	pass := false

	nb_styles := 0
	for _, line := range p.LyricsFile.Styles {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				nb_styles += 1
				break
			}
		}
	}
	if nb_styles > 1 {
		// for the moment, we focus on single style karaoke
		pass = true
	} else {
		for _, line := range p.LyricsFile.Styles {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
					s, err := ParseStyle(strings.TrimPrefix(line, "Style: "))
					if err != nil {
						return nil, err
					}
					if s.SecondaryColour == "&H00FFFFFF" {
						pass = true
					}
					break
				}
			}
		}
	}
	report.Content["pass"] = strconv.FormatBool(pass)
	return report, nil
}
