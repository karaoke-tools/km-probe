// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/misc"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type StyleSingleWhite struct {
	baseprobe.BaseProbe
}

func NewStyleSingleWhite(karaData *karadata.KaraData) probe.Probe {
	return &StyleSingleWhite{
		baseprobe.New("style-single-white",
			"unfilled color is not white (only if single style)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Misc,
					Tags:    []uuid.UUID{misc.GroupSinging},
					Msg:     "group singing karaoke", // we can use one color by voice
				},
				cond.HasMoreTagsThan{
					TagType: tag.Langs,
					Number:  1,
					Msg:     "is multilingual karaoke", // we can use one color by language
				},
				cond.HasAnyTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.MUL},
					Msg:     "is multilingual karaoke", // we can use one color by language
				},
			},
			karaData),
	}
}

func (p *StyleSingleWhite) Run(ctx context.Context) (report.Report, error) {
	nb_styles := 0
	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				nb_styles += 1
				if nb_styles > 1 {
					// for the moment, we focus on single style karaoke
					return report.Fail(severity.Warning, "multiple styles: check if this is a group singing karaoke (and add the tag if it is), or invert choir style colors"), nil
				}
			}
		}
	}
	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if s.SecondaryColour == "&H00FFFFFF" {
					// secondary color must be white if single style karaoke
					return report.Pass(), nil
				}
				break
			}
		}
	}
	return report.Fail(severity.Critical, "update style: secondary color must be white"), nil
}
