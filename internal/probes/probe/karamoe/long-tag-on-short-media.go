// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/misc"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"
)

type LongTagOnShortMedia struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewLongTagOnShortMedia() probe.Probe {
	return &LongTagOnShortMedia{
		baseprobe.New("long-tag-on-short-media",
			"long tag added manually",
			cond.GreaterMediaDuration{Duration: 300},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p LongTagOnShortMedia) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Misc, misc.Long) {
		return report.Fail(severity.Critical, "remove long tag"), nil
	}
	return report.Pass(), nil
}
