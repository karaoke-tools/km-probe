// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/version"
)

type VersionConflict struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewVersionConflict() lint.Lint {
	return &VersionConflict{
		baselint.New(
			"version-conflict",
			"incompatible version tags",
			cond.HasLessTagsThan{
				TagType: tag.Versions,
				Number:  2,
				Msg:     "has not multiple version tags",
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p VersionConflict) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Short) &&
		slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Full) {
		return report.Fail(severity.Critical, "song cannot be both a short version and a full version at the same time"), nil
	}
	if slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Acoustic) &&
		slices.Contains(KaraData.KaraJson.Data.Tags.Versions, version.Metal) {
		return report.Fail(severity.Critical, "song cannot be both a metal version and an acoustic version at the same time"), nil
	}
	return report.Pass(), nil
}
