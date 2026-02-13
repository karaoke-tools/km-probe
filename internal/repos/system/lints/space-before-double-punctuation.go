// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"

	"github.com/gofrs/uuid/v5"
)

type SpaceBeforeDoublePunctuation struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewSpaceBeforeDoublePunctuation() lint.Lint {
	return &SpaceBeforeDoublePunctuation{
		baselint.New("space-before-double-punctuation",
			"space before double punctuation (JPN/ENG only)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasTagsNotFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN, language.ENG},
					Msg:     "non english/japanese language",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p SpaceBeforeDoublePunctuation) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
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
