// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/misc"
	"github.com/louisroyer/km-probe/internal/karajson/songtype"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type MusicVideoCreditless struct {
	baseProbe
}

func NewMusicVideoCreditless(karaData *karadata.KaraData) Probe {
	return &MusicVideoCreditless{
		newBaseProbe("music-video-creditless", karaData),
	}
}

func (p *MusicVideoCreditless) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.KaraJson.Data.Tags.Misc) == 0 {
		return report.Skip("no misc tag"), nil
	}
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, songtype.MusicVideo) {
		return report.Skip("not a music video"), nil
	}

	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, misc.Creditless) {
		return report.Fail(severity.Critical, "music videos cannot be creditless, remove this tag"), nil
	}

	return report.Pass(), nil
}
