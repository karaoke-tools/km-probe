// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type ScaledBorderAndShadow struct {
	baseProbe
}

func NewScaledBorderAndShadow(karaData *karadata.KaraData) Probe {
	return &ScaledBorderAndShadow{
		newBaseProbe("scaled-border-and-shadow", karaData),
	}
}

func (p *ScaledBorderAndShadow) Run(ctx context.Context) (report.Report, error) {
	if p.karaData.Lyrics.ScriptInfo.ScaledBorderAndShadow {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "check the \"Scale border and shadow\" box"), nil
}
