// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"github.com/gofrs/uuid"

	"github.com/louisroyer/km-probe/internal/karajson"
)

// State of "no live download" collections as of 2025-01-06
var collectionsNoLiveDownload = []uuid.UUID{
	karajson.CollectionAsia,
	karajson.CollectionKana,
	karajson.CollectionWest,
}

var checkLiveDownloadProbablyAllowedKey = "live-download-probably-allowed"

// Checking each tag may be long when probing the full repository.
// This function only check for hardcoded collections and "unavailable" tag.
func (p *Probe) checkLiveDownloadProbablyAllowed(ctx context.Context) error {
	for _, tag := range p.KaraJson.Data.Tags.Misc {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if tag == karajson.Unavailable {
				p.Report.Fail(checkLiveDownloadProbablyAllowedKey)
				return nil
			}
		}
	}
	for _, tag := range p.KaraJson.Data.Tags.Collections {
		for _, collection := range collectionsNoLiveDownload {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if tag == collection {
					p.Report.Fail(checkLiveDownloadProbablyAllowedKey)
					return nil
				}
			}
		}
	}
	p.Report.Pass(checkLiveDownloadProbablyAllowedKey)
	return nil
}
