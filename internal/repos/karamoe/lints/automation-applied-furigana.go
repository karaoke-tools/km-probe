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

type AutomationAppliedFurigana struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewAutomationAppliedFurigana() lint.Lint {
	return &AutomationAppliedFurigana{
		baselint.New("automation-applied-furigana",
			"automation script not applied (song with furigana) ",
			cond.Any{
				cond.NoLyrics{},
				// we skip this lint when this is not a song in kana
				cond.HasNoTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.NonLatin},
					Msg:     "song in latin script",
				},
				cond.HasNoTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN},
					Msg:     "not a japanese song",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

// This is a generic version of "`automation-applied` lint" where we only check if at least one
// line has been generated from automation script and no line with "karaoke" effect is uncommented.
func (p AutomationAppliedFurigana) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	fx := false
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if line.Type == lyrics.Dialogue {
				switch line.Effect {
				case "fx":
					fx = true
				case "karaoke":
					return report.Fail(severity.Critical, "automation script has not been applied"), nil
				}
			}
		}
	}
	if fx {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "automation script has not been applied"), nil
}
