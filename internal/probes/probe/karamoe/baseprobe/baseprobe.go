// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

import (
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/skip"
)

type BaseProbe struct {
	baseprobe.BaseProbe
}

func New(name string, desc string, skipCond skip.Condition, karaData *karadata.KaraData) BaseProbe {
	return BaseProbe{
		baseprobe.New("kara-moe", name, desc, skipCond, karaData),
	}
}
