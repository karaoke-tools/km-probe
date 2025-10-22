// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/misc"
	"github.com/karaoke-tools/km-probe/internal/karajson/system/songtype"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid/v5"
)

type AudioOnlyCreditless struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewAudioOnlyCreditless() probe.Probe {
	return &AudioOnlyCreditless{
		baseprobe.New(
			"audio-only-creditless",
			"audio only songs cannot be creditless",
			cond.Any{
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.AUDIO},
					Msg:     "song is not audio only",
				},
				cond.HasNoTagFrom{
					TagType: tag.Misc,
					Tags:    []uuid.UUID{misc.Creditless},
					Msg:     "song is not tagged as creditless",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p AudioOnlyCreditless) Run(ctx context.Context, karaData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "remove the creditless tag, or update the media content"), nil
}
