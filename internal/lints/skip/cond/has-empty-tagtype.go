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

type HasEmptyTagtype struct {
	TagType tag.Tag
	Msg     string
}

func (h HasEmptyTagtype) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if k.KaraJson.HasEmptyTagtype(h.TagType) {
		return true, h.Msg, nil
	}
	return false, "", nil
}
