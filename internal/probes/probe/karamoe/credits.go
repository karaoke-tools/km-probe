// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/misc"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/origin"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/songtype"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type Credits struct {
	baseprobe.BaseProbe
}

var possibleCreditlessOrigin = []uuid.UUID{
	origin.Movie,
	origin.OriginalNetworkAnimation,
	origin.OriginalVideoAnimation,
	origin.TVSpecial,
	origin.TVSeries,
}

func NewCredits() probe.Probe {
	return &Credits{
		baseprobe.New(
			"credits",
			"can a creditless version be found?",
			cond.Any{
				cond.HasAnyTagFrom{
					TagType: tag.Misc,
					Tags:    []uuid.UUID{misc.Creditless},
					Msg:     "is a creditless version",
				},
				cond.HasAnyTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.AUDIO},
					Msg:     "has audio only tag",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED},
					Msg:     "songtype is not OP/ED",
				},
				cond.HasNoTagFrom{
					TagType: tag.Origins,
					Tags:    possibleCreditlessOrigin,
					Msg:     "origin not compatible with creditless",
				},
			},
		),
	}
}

func (p Credits) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Info, "if the media is already creditless, add the `Creditless`; if a creditless version exists (and is relevant!! see <https://kara.moe/playlist/quand-le-staff-fait-parti-du-generique> for counter-examples), update the media and add the tag"), nil
}
