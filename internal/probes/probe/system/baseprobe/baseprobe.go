// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

import (
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe/baseprobe"
)

type BaseProbe struct {
	baseprobe.BaseProbe
}

func New(name string, desc string, karaData *karadata.KaraData) BaseProbe {
	return BaseProbe{
		baseprobe.New("system", name, desc, karaData),
	}
}
