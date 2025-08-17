// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type NoLyrics struct {
	baseprobe.BaseProbe
}

func NewNoLyrics(karaData *karadata.KaraData) probe.Probe {
	return &NoLyrics{
		baseprobe.New("no-lyrics",
			"missing lyrics file",
			karaData),
	}
}

func (p *NoLyrics) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.Lyrics) > 0 {
		return report.Pass(), nil // contains a lyrics file
	}
	if res, err := p.KaraData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.ZXX}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Fail(severity.Critical, "no lyrics file, but the media is supposed to have has linguistic content"), nil
	}
	return report.Pass(), nil // no linguistical content
}
