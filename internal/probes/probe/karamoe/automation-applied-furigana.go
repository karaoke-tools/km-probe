// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/collection"
	"github.com/karaoke-tools/km-probe/internal/karajson/system/language"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid/v5"
)

type AutomationAppliedFurigana struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewAutomationAppliedFurigana() probe.Probe {
	return &AutomationAppliedFurigana{
		baseprobe.New("automation-applied-furigana",
			"automation script not applied (karaokes with furigana) ",
			cond.Any{
				cond.NoLyrics{},
				// we skip this probe when this is not a karaoke in kana
				cond.HasNoTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.NonLatin},
					Msg:     "karaoke in latin script",
				},
				cond.HasNoTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN},
					Msg:     "not a japanese karaoke",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

// This is a generic version of "`automation-applied` probe" where we only check if at least one
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
