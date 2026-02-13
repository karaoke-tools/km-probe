// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
)

type HasNoSongorder struct{}

func (h HasNoSongorder) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if k.KaraJson.Data.Songorder == nil {
		return true, "has no songorder", nil
	}
	return false, "", nil
}
