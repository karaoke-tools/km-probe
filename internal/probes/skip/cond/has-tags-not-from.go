// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/tag"

	"github.com/gofrs/uuid"
)

type HasTagsNotFrom struct {
	TagType tag.Tag
	Tags    []uuid.UUID
	Msg     string
}

func (h HasTagsNotFrom) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	if b, err := k.KaraJson.HasOnlyTagsFrom(ctx, h.TagType, h.Tags); err != nil {
		return true, "", err
	} else if !b {
		return true, h.Msg, nil
	}
	return false, "", nil
}
