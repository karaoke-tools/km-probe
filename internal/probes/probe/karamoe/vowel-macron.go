// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/collection"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type VowelMacron struct {
	baseprobe.BaseProbe
}

func NewVowelMacron(karaData *karadata.KaraData) probe.Probe {
	return &VowelMacron{
		baseprobe.New("vowel-macron",
			"ā, ē, ō, ī, ū in lyrics file (JPN romaji only)",
			karaData),
	}
}

func (p *VowelMacron) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	// we only check if language is full jpn romaji
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Collections, collection.Kana) {
		// TODO: with multi-track drifting, maybe we will have a way to detect kana versions
		// with System tags? so we could put this as system probe…
		return report.Skip("kana karaoke"), nil
	}
	if res, err := p.KaraData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip("not a japanese karaoke"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Events {
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
