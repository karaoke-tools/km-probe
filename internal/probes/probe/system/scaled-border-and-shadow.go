// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type ScaledBorderAndShadow struct {
	baseprobe.BaseProbe
}

func NewScaledBorderAndShadow(karaData *karadata.KaraData) probe.Probe {
	return &ScaledBorderAndShadow{
		baseprobe.New("scaled-border-and-shadow",
			"scaled border and shadow not enabled",
			karaData),
	}
}

func (p *ScaledBorderAndShadow) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}

	// TODO: update this when multi-track drifting is released
	if p.KaraData.Lyrics[0].ScriptInfo.ScaledBorderAndShadow {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "check the \"Scale border and shadow\" box"), nil
}
