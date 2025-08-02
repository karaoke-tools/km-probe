// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/version"
	"github.com/louisroyer/km-probe/internal/probes/report"

	"github.com/gofrs/uuid"
)

type AltVersionWithoutParent struct {
	baseProbe
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

func NewAltVersionWithoutParent(karaData *karadata.KaraData) Probe {
	return &AltVersionWithoutParent{
		newBaseProbe("alt-version-without-parent", karaData),
	}
}

func (p *AltVersionWithoutParent) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.KaraJson.Data.Parents) > 0 {
		return report.Skip(), nil
	}
	if len(p.karaData.KaraJson.Data.Tags.Versions) == 0 {
		return report.Skip(), nil
	}
	if slices.ContainsFunc(p.karaData.KaraJson.Data.Tags.Versions, isVersionWithoutParentCritical) {
		return report.Fail(), nil
	}
	return report.Pass(), nil
}
