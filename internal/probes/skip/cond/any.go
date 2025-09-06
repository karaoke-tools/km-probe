// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/probes/skip"
)

type Any []skip.Condition

func (a Any) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	for _, c := range a {
		if ok, msg, err := c.Result(ctx, k); err != nil {
			return true, "", err
		} else if ok {
			return ok, msg, nil
		}
	}
	return false, "", nil
}
