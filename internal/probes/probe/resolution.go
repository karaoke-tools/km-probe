// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type Resolution struct {
	baseProbe
}

func NewResolution(karaData *karadata.KaraData) Probe {
	return &Resolution{
		newBaseProbe("resolution", karaData),
	}
}

func (p *Resolution) Run(ctx context.Context) (report.Report, error) {
	if p.karaData.Lyrics.ScriptInfo.PlayResX == 0 && p.karaData.Lyrics.ScriptInfo.PlayResY == 0 {
		return report.Pass(), nil
	}
	return report.Fail(), nil
}
