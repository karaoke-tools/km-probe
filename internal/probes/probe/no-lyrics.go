// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/language"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type NoLyrics struct {
	baseProbe
}

func NewNoLyrics(karaData *karadata.KaraData) Probe {
	return &NoLyrics{
		newBaseProbe("no-lyrics", karaData),
	}
}

func (p *NoLyrics) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.Lyrics) > 0 {
		return report.Pass(), nil // contains a lyrics file
	}
	if res, err := p.karaData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.ZXX}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Fail(severity.Critical, "no lyrics file, but the media is supposed to have has linguistic content"), nil
	}
	return report.Pass(), nil // no linguistical content
}
