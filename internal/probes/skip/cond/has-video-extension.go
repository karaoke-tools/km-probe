// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
)

type HasVideoExtension struct{}

func (h HasVideoExtension) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	// TODO: multi-track drifting
	filename := k.KaraJson.Medias[0].Filename
	startExt := strings.LastIndexByte(filename, '.')
	extension := filename[startExt+1:]
	if !slices.Contains([]string{"mp4", "mkv", "webm"}, extension) {
		return true, "has video extension", nil
	}
	return false, "", nil
}
