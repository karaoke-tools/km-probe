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
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
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
	if len(p.karaData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, misc.GroupSinging) {
		return report.Skip("group singing karaoke: secondary color can be non white"), nil
	}

	nb_styles := 0
	// TODO: update this when multi-track drifting is released
	for _, line := range p.karaData.Lyrics[0].Styles {
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
	for _, line := range p.karaData.Lyrics[0].Styles {
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
