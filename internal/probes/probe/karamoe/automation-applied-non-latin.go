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

type AutomationAppliedNonLatin struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewAutomationAppliedNonLatin() probe.Probe {
	return &AutomationAppliedNonLatin{
		baseprobe.New("automation-applied-non-latin",
			"automation script not applied (non-latin scripts only) ",
			cond.Any{
				cond.NoLyrics{},
				cond.HasNoTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.Kana},
					Msg:     "karaoke in latin script",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

// This is a generic version of "`automation-applied` probe" where we only check if at least one
// line has been generated from automation script and no line with "karaoke" effect is uncommented.
func (p AutomationAppliedNonLatin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
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
