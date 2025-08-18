// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karajson

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karajson/tag"

	"github.com/gofrs/uuid"
)

func (k KaraJson) Tag(t tag.Tag) []uuid.UUID {
	switch t {
	case tag.Authors:
		return k.Data.Tags.Authors
	case tag.Collections:
		return k.Data.Tags.Collections
	case tag.Creators:
		return k.Data.Tags.Creators
	case tag.Families:
		return k.Data.Tags.Families
	case tag.Groups:
		return k.Data.Tags.Groups
	case tag.Langs:
		return k.Data.Tags.Langs
	case tag.Misc:
		return k.Data.Tags.Misc
	case tag.Origins:
		return k.Data.Tags.Origins
	case tag.Platforms:
		return k.Data.Tags.Platforms
	case tag.Series:
		return k.Data.Tags.Series
	case tag.Singers:
		return k.Data.Tags.Singers
	case tag.Singergroups:
		return k.Data.Tags.Singergroups
	case tag.Songtypes:
		return k.Data.Tags.Songtypes
	case tag.Songwriters:
		return k.Data.Tags.Songwriters
	case tag.Versions:
		return k.Data.Tags.Versions
	case tag.Warnings:
		return k.Data.Tags.Warnings
	}
	return []uuid.UUID{}
}

func (k *KaraJson) HasOnlyTagsFrom(ctx context.Context, t tag.Tag, ids []uuid.UUID) (bool, error) {
	field := k.Tag(t)
	if len(field) > len(ids) {
		return false, nil
	}
	for _, l := range field {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			if !slices.Contains(ids, l) {
				return false, nil
			}
		}
	}
	return true, nil
}

func (k *KaraJson) HasAnyTagFrom(t tag.Tag, ids []uuid.UUID) bool {
	field := k.Tag(t)
	return slices.ContainsFunc(field, func(u uuid.UUID) bool { return slices.Contains(ids, u) })
}

func (k *KaraJson) HasEmptyTagtype(t tag.Tag) bool {
	return len(k.Tag(t)) == 0
}
