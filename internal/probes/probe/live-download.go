// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"

	"github.com/gofrs/uuid"
)

type LiveDownload struct {
	baseProbe
}

func NewLiveDownload(karaData *karadata.KaraData) Probe {
	return &LiveDownload{
		newBaseProbe("live-download", karaData),
	}
}

// State of "no live download" collections as of 2025-01-06
var collectionsNoLiveDownload = []uuid.UUID{
	karajson.CollectionAsia,
	karajson.CollectionKana,
	karajson.CollectionWest,
}

func isNoLiveDownloadCollection(collection uuid.UUID) bool {
	return slices.Contains(collectionsNoLiveDownload, collection)
}

// Checking each tag may be long when probing the full repository.
// This function only check for hardcoded collections and "unavailable" tag.
func (p *LiveDownload) Run(ctx context.Context) (report.Report, error) {
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, karajson.Unavailable) {
		return report.Info(false), nil
	}
	if slices.ContainsFunc(p.karaData.KaraJson.Data.Tags.Collections, isNoLiveDownloadCollection) {
		return report.Info(false), nil
	}
	return report.Info(true), nil
}
