// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
)

type NoLyrics struct{}

func (n NoLyrics) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if len(k.Lyrics) == 0 {
		return true, "no lyrics", nil
	}
	return false, "", nil
}
