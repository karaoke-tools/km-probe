// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/probes/skip"
)

type BaseProbe struct {
	pkg      string
	name     string
	desc     string
	skipCond skip.Condition
}

func New(pkg string, name string, desc string, skipCond skip.Condition) BaseProbe {
	return BaseProbe{
		name:     strings.Join([]string{pkg, name}, "."),
		desc:     desc,
		skipCond: skipCond,
	}
}

func (p BaseProbe) Name() string {
	return p.name
}

func (p BaseProbe) Description() string {
	return p.desc
}

func (p BaseProbe) PreRun(ctx context.Context, KaraData *karadata.KaraData) (bool, string, error) {
	return p.skipCond.Result(ctx, KaraData)
}
