// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/collection"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type VowelMacron struct {
	baseprobe.BaseProbe
}

func NewVowelMacron() probe.Probe {
	return &VowelMacron{
		baseprobe.New("vowel-macron",
			"ā, ē, ō, ī, ū in lyrics file (JPN romaji only)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.Kana},
					Msg:     "kana karaoke",
				},
				cond.HasTagsNotFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN},
					Msg:     "not a japanese only karaoke",
				},
			},
		),
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
					return report.Fail(severity.Warning, "in full japanese karaoke, vowels should not have macron: use the appropriate expansion from: aa/ii/uu/ee/ou/oo (if this is on a chinese word, make sure to put it in fullcaps)"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
