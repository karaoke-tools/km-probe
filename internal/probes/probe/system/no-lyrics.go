// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type NoLyrics struct {
	baseprobe.BaseProbe
}

func NewNoLyrics() probe.Probe {
	return &NoLyrics{
		baseprobe.New("no-lyrics",
			"missing lyrics file",
			cond.HasLyrics{},
		),
	}
}

func (p NoLyrics) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if res := KaraData.KaraJson.HasAnyTagFrom(tag.Langs, []uuid.UUID{language.ZXX}); !res {
		return report.Fail(severity.Critical, "no lyrics file, but the media is supposed to have has linguistic content"), nil
	}
	return report.Pass(), nil // no linguistical content
}
