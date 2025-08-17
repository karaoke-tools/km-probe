// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/collection"
	"github.com/louisroyer/km-probe/internal/karajson/language"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type VowelMacron struct {
	baseProbe
}

func NewVowelMacron(karaData *karadata.KaraData) Probe {
	return &VowelMacron{
		newBaseProbe("vowel-macron", karaData),
	}
}

func (p *VowelMacron) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	// we only check if language is full jpn romaji
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Collections, collection.Kana) {
		return report.Skip("kana karaoke"), nil
	}
	if res, err := p.karaData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip("not a japanese karaoke"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.karaData.Lyrics[0].Events {
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
