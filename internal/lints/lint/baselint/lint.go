// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baselint

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/skip"
)

type BaseLint struct {
	pkg      string
	name     string
	desc     string
	skipCond skip.Condition
}

func New(pkg string, name string, desc string, skipCond skip.Condition) BaseLint {
	return BaseLint{
		pkg:      pkg,
		name:     name,
		desc:     desc,
		skipCond: skipCond,
	}
}

func (p BaseLint) Name() string {
	return strings.Join([]string{p.pkg, p.name}, ".")
}

func (p BaseLint) Description() string {
	return p.desc
}

func (p BaseLint) PreRun(ctx context.Context, KaraData *karadata.KaraData) (bool, string, error) {
	return p.skipCond.Result(ctx, KaraData)
}
