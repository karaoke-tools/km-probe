// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
)

type Never struct{}

func (n Never) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	return false, "", nil
}
