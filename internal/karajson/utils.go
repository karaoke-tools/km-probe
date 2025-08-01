// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karajson

import (
	"context"
	"slices"

	"github.com/gofrs/uuid"
)

// Ensure the karajson includes only languages from the given slice.
func (k KaraJson) HasOnlyLanguagesFrom(ctx context.Context, langs []uuid.UUID) (bool, error) {
	if len(k.Data.Tags.Langs) > len(langs) {
		return false, nil
	}
	for _, l := range k.Data.Tags.Langs {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			if !slices.Contains(langs, l) {
				return false, nil
			}
		}
	}
	return true, nil
}
