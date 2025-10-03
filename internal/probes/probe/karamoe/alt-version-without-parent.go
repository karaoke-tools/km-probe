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

	"github.com/gofrs/uuid"
)

type AltVersionWithoutParent struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

var versionsWithoutParentInfo = []uuid.UUID{
	version.Cover,
	version.OffVocal,
	version.NonLatin,
}

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

func isVersionWithoutParentInfo(versionType uuid.UUID) bool {
	return slices.Contains(versionsWithoutParentInfo, versionType)
}

func NewAltVersionWithoutParent() probe.Probe {
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
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p AltVersionWithoutParent) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.ContainsFunc(KaraData.KaraJson.Data.Tags.Versions, isVersionWithoutParentCritical) {
		return report.Fail(severity.Critical, "check if a potential parent exists, or if the version tag is relevant"), nil
	}

	if slices.ContainsFunc(KaraData.KaraJson.Data.Tags.Versions, isVersionWithoutParentInfo) {
		return report.Fail(severity.Info, "has a version tag but no parent: if the parent exist, don't forget to add it"), nil
	}
	return report.Pass(), nil
}
