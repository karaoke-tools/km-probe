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
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type LongTagOnShortMedia struct {
	baseProbe
}

func NewLongTagOnShortMedia(karaData *karadata.KaraData) Probe {
	return &LongTagOnShortMedia{
		newBaseProbe("long-tag-on-short-media", karaData),
	}
}

func (p *LongTagOnShortMedia) Run(ctx context.Context) (report.Report, error) {
	if p.karaData.KaraJson.Medias[0].Duration > 300 {
		return report.Skip("media over 300s"), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, misc.Long) {
		return report.Fail(severity.Critical, "remove long tag"), nil
	}
	return report.Pass(), nil
}
