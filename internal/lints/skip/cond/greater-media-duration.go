// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"
	"fmt"

	"github.com/karaoke-tools/km-probe/internal/karadata"
)

type GreaterMediaDuration struct {
	Duration int
}

func (g GreaterMediaDuration) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if k.KaraJson.Medias[0].Duration > 300 {
		return true, fmt.Sprintf("media duration over %d seconds", g.Duration), nil
	}
	return false, "", nil
}
