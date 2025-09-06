// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

import (
	"github.com/karaoke-tools/km-probe/internal/probes/probe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/skip"
)

type BaseProbe struct {
	baseprobe.BaseProbe
}

func New(name string, desc string, skipCond skip.Condition) BaseProbe {
	return BaseProbe{
		baseprobe.New("kara-moe", name, desc, skipCond),
	}
}
