// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/language"
	"github.com/louisroyer/km-probe/internal/karajson/warning"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type LyricsWarningNoLinguisticalContent struct {
	baseProbe
}

func NewLyricsWarningNoLinguisticalContent(karaData *karadata.KaraData) Probe {
	return &LyricsWarningNoLinguisticalContent{
		newBaseProbe("lyrics-warning-no-linguistical-content", karaData),
	}
}

func (p *LyricsWarningNoLinguisticalContent) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Warnings, warning.R18Lyrics) {
		return report.Skip("no lyrics warning"), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Langs, language.ZXX) {
		return report.Fail(severity.Critical, "check if lyrics warning is relevant, and if the Langs field is set"), nil
	}
	return report.Pass(), nil
}
