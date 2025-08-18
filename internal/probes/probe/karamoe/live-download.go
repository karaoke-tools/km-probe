// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/collection"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/misc"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type LiveDownload struct {
	baseprobe.BaseProbe
}

func NewLiveDownload(karaData *karadata.KaraData) probe.Probe {
	return &LiveDownload{
		baseprobe.New(
			"live-download",
			"is hardsub available?",
			cond.Never{},
			karaData),
	}
}

// State of "no live download" collections as of 2025-01-06
var collectionsNoLiveDownload = []uuid.UUID{
	collection.Asia,
	collection.Kana,
	collection.West,
}

func isNoLiveDownloadCollection(collection uuid.UUID) bool {
	return slices.Contains(collectionsNoLiveDownload, collection)
}

// Checking each tag may be long when probing the full repository.
// This function only check for hardcoded collections and "unavailable" tag.
func (p *LiveDownload) Run(ctx context.Context) (report.Report, error) {
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Misc, misc.Unavailable) {
		return report.Info(false), nil
	}
	if slices.ContainsFunc(p.KaraData.KaraJson.Data.Tags.Collections, isNoLiveDownloadCollection) {
		return report.Info(false), nil
	}
	return report.Info(true), nil
}
