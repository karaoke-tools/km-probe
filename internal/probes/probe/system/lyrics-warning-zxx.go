// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/karajson/system/warning"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type LyricsWarningZXX struct {
	baseprobe.BaseProbe
}

func NewLyricsWarningZXX(karaData *karadata.KaraData) probe.Probe {
	return &LyricsWarningZXX{
		baseprobe.New("lyrics-warning-zxx",
			"lyrics warning, but there is no linguistical content",
			karaData),
	}
}

func (p *LyricsWarningZXX) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.KaraData.KaraJson.Data.Tags.Warnings, warning.R18Lyrics) {
		return report.Skip("no lyrics warning"), nil
	}
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Langs, language.ZXX) {
		return report.Fail(severity.Critical, "check if lyrics warning is relevant, and if the Langs field is set"), nil
	}
	return report.Pass(), nil
}
