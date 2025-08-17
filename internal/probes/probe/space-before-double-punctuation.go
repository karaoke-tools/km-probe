// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/language"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type SpaceBeforeDoublePunctuation struct {
	baseProbe
}

func NewSpaceBeforeDoublePunctuation(karaData *karadata.KaraData) Probe {
	return &SpaceBeforeDoublePunctuation{
		newBaseProbe("space-before-double-punctuation", karaData),
	}
}

func (p *SpaceBeforeDoublePunctuation) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	// we only check if language is full english, full japanese, or jpn+eng
	if res, err := p.karaData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN, language.ENG}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip("non english/japanese language"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.karaData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && !((line.Type == lyrics.Comment) && (line.Effect != "karaoke")) {
				l := line.Text.StripTags()
				if strings.Contains(l, " ?") || strings.Contains(l, " !") ||
					strings.Contains(l, " ?") || strings.Contains(l, " !") { // non-breakable space
					return report.Fail(severity.Critical, "remove space before `?`/`!`"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
