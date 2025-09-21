// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/version"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"
)

type VersionConflict struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewVersionConflict() probe.Probe {
	return &VersionConflict{
		baseprobe.New(
			"version-conflict",
			"incompatible version tags",
			cond.HasLessTagsThan{
				TagType: tag.Versions,
				Number:  2,
				Msg:     "has not multiple version tags",
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p VersionConflict) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Short) &&
		slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Full) {
		return report.Fail(severity.Critical, "karaoke cannot be both a short version and a full version at the same time"), nil
	}
	if slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Acoustic) &&
		slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Metal) {
		return report.Fail(severity.Critical, "karaoke cannot be both a metal version and an acoustic version at the same time"), nil
	}
	return report.Pass(), nil
}
