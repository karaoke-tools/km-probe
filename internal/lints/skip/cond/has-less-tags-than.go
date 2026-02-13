// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
)

type HasLessTagsThan struct {
	TagType tag.Tag
	Number  int
	Msg     string
}

func (h HasLessTagsThan) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if len(k.KaraJson.Tag(h.TagType)) < h.Number {
		return true, h.Msg, nil
	}
	return false, "", nil
}
