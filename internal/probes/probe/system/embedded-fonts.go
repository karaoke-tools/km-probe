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

type EmbeddedFonts struct {
	baseprobe.BaseProbe
}

func NewEmbeddedFonts() probe.Probe {
	return &EmbeddedFonts{
		baseprobe.New("embedded-fonts",
			"lyrics file embeds fonts",
			cond.NoLyrics{},
		),
	}
}

func (p EmbeddedFonts) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].Fonts {
		return report.Fail(severity.Critical, "lyrics file embeds fonts; consider using standard fonts instead because fonts embedding creates big lyrics file"), nil
	}
	return report.Pass(), nil
}
