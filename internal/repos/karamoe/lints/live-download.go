// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/collection"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/misc"

	"github.com/gofrs/uuid/v5"
)

type LiveDownload struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewLiveDownload() lint.Lint {
	return &LiveDownload{
		baselint.New(
			"live-download",
			"is hardsub available?",
			cond.Never{},
		),
		baselint.EnabledByDefault{},
	}
}

// State of "no live download" collections as of 2025-01-06
var collectionsNoLiveDownload = []uuid.UUID{
	collection.Asia,
	collection.NonLatin,
	collection.West,
}

func isNoLiveDownloadCollection(collection uuid.UUID) bool {
	return slices.Contains(collectionsNoLiveDownload, collection)
}

// Checking each tag may be long when probing the full repository.
// This function only check for hardcoded collections and "unavailable" tag.
func (p LiveDownload) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Misc, misc.Unavailable) {
		return report.Fail(severity.Info, "not available for live download"), nil
	}
	if slices.ContainsFunc(KaraData.KaraJson.Data.Tags.Collections, isNoLiveDownloadCollection) {
		return report.Fail(severity.Info, "not available for live download"), nil
	}
	return report.Pass(), nil
}
