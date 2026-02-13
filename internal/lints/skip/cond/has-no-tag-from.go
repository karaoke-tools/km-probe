// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"

	"github.com/gofrs/uuid/v5"
)

type HasNoTagFrom struct {
	TagType tag.Tag
	Tags    []uuid.UUID
	Msg     string
}

func (h HasNoTagFrom) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if b := k.KaraJson.HasAnyTagFrom(h.TagType, h.Tags); !b {
		return true, h.Msg, nil
	}
	return false, "", nil
}
