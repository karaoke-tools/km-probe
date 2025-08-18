// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cond

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/skip"
)

type All []skip.Condition

func (a All) Result(ctx context.Context, k *karadata.KaraData) (bool, string, error) {
	msgs := []string{}
	for _, c := range a {
		if ok, msg, err := c.Result(ctx, k); err != nil {
			return true, "", err
		} else if !ok {
			return ok, msg, nil
		} else {
			msgs = append(msgs, msg)
		}
	}
	return true, strings.Join(msgs, "; "), nil
}
