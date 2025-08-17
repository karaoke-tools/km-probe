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
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type MusicVideoCreditless struct {
	baseprobe.BaseProbe
}

func NewMusicVideoCreditless(karaData *karadata.KaraData) probe.Probe {
	return &MusicVideoCreditless{
		baseprobe.New(
			"music-video-creditless",
			"MV with a creditless tag",
			karaData),
	}
}

func (p *MusicVideoCreditless) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.KaraJson.Data.Tags.Misc) == 0 {
		return report.Skip("no misc tag"), nil
	}
	if !slices.Contains(p.KaraData.KaraJson.Data.Tags.Songtypes, songtype.MusicVideo) {
		return report.Skip("not a music video"), nil
	}

	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Misc, misc.Creditless) {
		return report.Fail(severity.Critical, "music videos cannot be creditless, remove this tag"), nil
	}

	return report.Pass(), nil
}
