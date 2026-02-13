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
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/collection"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"

	"github.com/gofrs/uuid/v5"
)

type WrongTsuSeparation struct {
	baselint.BaseLint
	lint.WithDefault
}

// In japanese, `つ` should be timed as a single syllable: `tsu`.
// For example, `ひとつ`(`hitotsu`) should be timed as `hi|to|tsu` and not as `hi|tot|su`.
func NewWrongTsuSeparation() lint.Lint {
	return &WrongTsuSeparation{
		baselint.New("wrong-tsu-separation",
			"`t|su` separation is not correct (JPN romaji only)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.NonLatin},
					Msg:     "non-latin script song",
				},
				cond.HasTagsNotFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN},
					Msg:     "not a japanese only song",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p WrongTsuSeparation) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				ok := false
				for _, syll := range line.Text.TagsSplit {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						if !strings.HasPrefix(syll, "{") {
							if strings.HasSuffix(syll, "t") {
								ok = true
							} else if ok && strings.HasPrefix(syll, "su") {
								return report.Fail(severity.Critical, "`tsu` must be timed as a single syllable"), nil
							} else {
								ok = false
							}
						}
					}
				}
			}
		}
	}
	return report.Pass(), nil
}
