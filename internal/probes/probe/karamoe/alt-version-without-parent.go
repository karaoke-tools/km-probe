// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/version"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type AltVersionWithoutParent struct {
	baseprobe.BaseProbe
}

// version.Cover, version.OffVocal and version.Kana are non critical
var versionsWithoutParentCritical = []uuid.UUID{
	version.Acoustic,
	version.Alternative,
	version.Full,
	version.Metal,
	version.Short,
}

func isVersionWithoutParentCritical(versionType uuid.UUID) bool {
	return slices.Contains(versionsWithoutParentCritical, versionType)
}

func NewAltVersionWithoutParent(karaData *karadata.KaraData) probe.Probe {
	return &AltVersionWithoutParent{
		baseprobe.New(
			"alt-version-without-parent",
			"version tag, but there is no parent song",
			cond.Any{
				cond.HasParent{},
				cond.HasEmptyTagtype{
					TagType: tag.Versions,
					Msg:     "not an alt version",
				},
			},
			karaData),
	}
}

func (p *AltVersionWithoutParent) Run(ctx context.Context) (report.Report, error) {
	if slices.ContainsFunc(p.KaraData.KaraJson.Data.Tags.Versions, isVersionWithoutParentCritical) {
		return report.Fail(severity.Critical, "check if a potential parent exists, or if the version tag is relevant"), nil
	}
	return report.Pass(), nil
}
