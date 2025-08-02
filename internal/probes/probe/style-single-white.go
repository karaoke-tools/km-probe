// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/misc"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type StyleSingleWhite struct {
	baseProbe
}

func NewStyleSingleWhite(karaData *karadata.KaraData) Probe {
	return &StyleSingleWhite{
		newBaseProbe("style-single-white", karaData),
	}
}

func (p *StyleSingleWhite) Run(ctx context.Context) (report.Report, error) {
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, misc.GroupSinging) {
		return report.Skip(), nil
	}
	nb_styles := 0
	for _, line := range p.karaData.Lyrics.Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				nb_styles += 1
				if nb_styles > 1 {
					// for the moment, we focus on single style karaoke
					return report.Skip(), nil
				}
			}
		}
	}
	for _, line := range p.karaData.Lyrics.Styles {
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
	return report.Fail(), nil
}
