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

type VowelMacron struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewVowelMacron() lint.Lint {
	return &VowelMacron{
		baselint.New("vowel-macron",
			"ā, ē, ō, ī, ū in lyrics file (JPN romaji only)",
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

func (p VowelMacron) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				if strings.ContainsAny(line.Text.StripTags(), "āīūēō") {
					return report.Fail(severity.Warning, "in full japanese song, vowels should not have macron: use the appropriate expansion from: aa/ii/uu/ee/ou/oo (if this is on a chinese word, make sure to put it in fullcaps)"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
