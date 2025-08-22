// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/misc"
	"github.com/louisroyer/km-probe/internal/karajson/system/songtype"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type AudioOnlyCreditless struct {
	baseprobe.BaseProbe
}

func NewAudioOnlyCreditless(karaData *karadata.KaraData) probe.Probe {
	return &AudioOnlyCreditless{
		baseprobe.New(
			"audio-only-creditless",
			"audo only songs cannot be creditless",
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
			karaData),
	}
}

func (p *AudioOnlyCreditless) Run(ctx context.Context) (report.Report, error) {
	return report.Fail(severity.Critical, "remove the creditless tag, or update the media content"), nil
}
