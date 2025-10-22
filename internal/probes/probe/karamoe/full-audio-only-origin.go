// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/origin"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/songtype"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/version"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid/v5"
)

type FullAudioOnlyOrigin struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewFullAudioOnlyOrigin() probe.Probe {
	return &FullAudioOnlyOrigin{
		baseprobe.New(
			"full-audio-only-origin",
			"audio only song cannot have origin when they have not the same size than actual OP/ED/IN",
			cond.Any{
				cond.HasEmptyTagtype{
					TagType: tag.Origins,
					Msg:     "has no origin",
				},
				cond.All{
					cond.HasLessTagsThan{
						TagType: tag.Origins,
						Number:  2,
						Msg:     "has a single origin tag",
					},
					cond.HasAnyTagFrom{
						TagType: tag.Origins,
						Tags:    []uuid.UUID{origin.Vtuber},
						Msg:     "vtuber tag",
					},
				},
				cond.HasNoTagFrom{
					TagType: tag.Versions,
					Tags:    []uuid.UUID{version.Full},
					Msg:     "is not a full version",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.AUDIO},
					Msg:     "is not an audio only",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p FullAudioOnlyOrigin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "remove origin tag"), nil
}
