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

type AegisubGarbage struct {
	baseprobe.BaseProbe
}

func NewAegisubGarbage() probe.Probe {
	return &AegisubGarbage{
		baseprobe.New("aegisub-garbage",
			"Aegisub Project Garbage section has not been removed",
			cond.NoLyrics{},
		),
	}
}

func (p AegisubGarbage) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].AegisubGarbage {
		return report.Fail(severity.Info, "Aegisub Project Garbage section has not been removed; if you are integrating this karaoke make sure to enable the \"cleanup lyrics\" function in Karaoke Mugen "), nil
	}
	return report.Pass(), nil
}
