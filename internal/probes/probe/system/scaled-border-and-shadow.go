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
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type ScaledBorderAndShadow struct {
	baseprobe.BaseProbe
}

func NewScaledBorderAndShadow() probe.Probe {
	return &ScaledBorderAndShadow{
		baseprobe.New("scaled-border-and-shadow",
			"scaled border and shadow not enabled",
			cond.NoLyrics{},
		),
	}
}

func (p ScaledBorderAndShadow) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].ScriptInfo.ScaledBorderAndShadow {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "check the \"Scale border and shadow\" box"), nil
}
