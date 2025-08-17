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
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type LongTagOnShortMedia struct {
	baseprobe.BaseProbe
}

func NewLongTagOnShortMedia(karaData *karadata.KaraData) probe.Probe {
	return &LongTagOnShortMedia{
		baseprobe.New("long-tag-on-short-media",
			"long tag added manually",
			karaData),
	}
}

func (p *LongTagOnShortMedia) Run(ctx context.Context) (report.Report, error) {
	if p.KaraData.KaraJson.Medias[0].Duration > 300 {
		return report.Skip("media over 300s"), nil
	}
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Misc, misc.Long) {
		return report.Fail(severity.Critical, "remove long tag"), nil
	}
	return report.Pass(), nil
}
