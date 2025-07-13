// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"github.com/gofrs/uuid"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"
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

// Checking each tag may be long when probing the full repository.
// This function only check for hardcoded collections and "unavailable" tag.
func (p *LiveDownload) Run(ctx context.Context) (report.Report, error) {
	for _, tag := range p.karaData.KaraJson.Data.Tags.Misc {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if tag == karajson.Unavailable {
				return report.Info(false), nil
			}
		}
	}
	for _, tag := range p.karaData.KaraJson.Data.Tags.Collections {
		for _, collection := range collectionsNoLiveDownload {
			select {
			case <-ctx.Done():
				return report.Abort(), ctx.Err()
			default:
				if tag == collection {
					return report.Info(false), nil
				}
			}
		}
	}
	return report.Info(true), nil
}
