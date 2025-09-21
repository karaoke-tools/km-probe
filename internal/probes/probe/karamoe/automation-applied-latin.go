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
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type AutomationAppliedLatin struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewAutomationAppliedLatin() probe.Probe {
	return &AutomationAppliedLatin{
		baseprobe.New("automation-applied-latin",
			"automation script not applied (latin scripts only) ",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.Kana},
					Msg:     "karaoke in non-latin script",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p AutomationAppliedLatin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
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
