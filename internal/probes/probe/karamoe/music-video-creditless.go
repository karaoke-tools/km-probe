// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

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

type MusicVideoCreditless struct {
	baseprobe.BaseProbe
}

func NewMusicVideoCreditless(karaData *karadata.KaraData) probe.Probe {
	return &MusicVideoCreditless{
		baseprobe.New(
			"music-video-creditless",
			"MV with a creditless tag",
			cond.Any{
				cond.NoLyrics{},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.MusicVideo},
					Msg:     "not a music video",
				},
			},
			karaData),
	}
}

func (p *MusicVideoCreditless) Run(ctx context.Context) (report.Report, error) {
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Misc, misc.Creditless) {
		return report.Fail(severity.Critical, "music videos cannot be creditless, remove this tag"), nil
	}

	return report.Pass(), nil
}
