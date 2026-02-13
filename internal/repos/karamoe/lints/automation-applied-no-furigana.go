// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/collection"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"

	"github.com/gofrs/uuid/v5"
)

type AutomationAppliedNoFurigana struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewAutomationAppliedNoFurigana() lint.Lint {
	return &AutomationAppliedNoFurigana{
		baselint.New("automation-applied-no-furigana",
			"automation script not applied (song without furigana) ",
			cond.Any{
				cond.NoLyrics{},
				cond.All{
					// we skip this lint when song is in furigana (non-latin with japanese)
					cond.HasAnyTagFrom{
						TagType: tag.Collections,
						Tags:    []uuid.UUID{collection.NonLatin},
						Msg:     "song with latin script",
					},
					cond.HasAnyTagFrom{
						TagType: tag.Langs,
						Tags:    []uuid.UUID{language.JPN},
						Msg:     "japanese song",
					},
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p AutomationAppliedNoFurigana) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	fx := 0
	karaoke := 0
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if line.Type == lyrics.Comment && line.Effect == "karaoke" {
				karaoke++
			} else if line.Type == lyrics.Dialogue {
				switch line.Effect {
				case "fx":
					fx++
				case "karaoke":
					return report.Fail(severity.Critical, "automation script has not been applied"), nil
				}
			}
		}
	}
	if fx == 0 || karaoke != fx {
		return report.Fail(severity.Critical, "automation script has not been applied"), nil
	}
	return report.Pass(), nil
}
