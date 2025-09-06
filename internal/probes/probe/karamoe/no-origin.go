// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/songtype"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type NoOrigin struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewNoOrigin() probe.Probe {
	return &NoOrigin{
		baseprobe.New(
			"no-origin",
			"songtype is OP/ED/IN but origin tag is missing",
			cond.Any{
				cond.HasMoreTagsThan{
					TagType: tag.Origins,
					Number:  0,
					Msg:     "has origin",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED, songtype.IN},
					Msg:     "songtype is not OP/ED/IN",
				},
				// for audio only OP/ED/IN, origin is not enforced:
				// motivation: origin tag is clearly relevant when it applies to both the media and the work (series tag),
				// but it is unclear if origin applies
				// - to the media (in this case, for audio-only it should generally not be filled), or
				// - to the work (in this case, why are origin tags not a property of the series tag (like database ids are)?)
				cond.HasAnyTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.AUDIO},
					Msg:     "is audio only",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p NoOrigin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "add the missing origin tag"), nil
}
