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

type Resolution struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewResolution() probe.Probe {
	return &Resolution{
		baseprobe.New("resolution",
			"resolution not set to 0×0",
			cond.NoLyrics{},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p Resolution) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].ScriptInfo.PlayResX == 0 && KaraData.Lyrics[0].ScriptInfo.PlayResY == 0 {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "update resolution to be 0×0 (and check style size)"), nil
}
