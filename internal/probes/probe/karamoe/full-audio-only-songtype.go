// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/songtype"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/version"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type FullAudioOnlySongtype struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewFullAudioOnlySongtype() probe.Probe {
	return &FullAudioOnlySongtype{
		baseprobe.New(
			"full-audio-only-songtype",
			"audio only song are considered OP/ED/IN only when they have the same size than actual OP/ED/IN",
			cond.Any{
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
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED, songtype.IN},
					Msg:     "is not an OP/ED/IN",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p FullAudioOnlySongtype) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "remove the OP/ED/IN tag"), nil
}
