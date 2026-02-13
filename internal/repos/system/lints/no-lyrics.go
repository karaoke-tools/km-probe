// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"

	"github.com/gofrs/uuid/v5"
)

type NoLyrics struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewNoLyrics() lint.Lint {
	return &NoLyrics{
		baselint.New("no-lyrics",
			"missing lyrics file",
			cond.HasLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p NoLyrics) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if res := KaraData.KaraJson.HasAnyTagFrom(tag.Langs, []uuid.UUID{language.ZXX}); !res {
		return report.Fail(severity.Critical, "no lyrics file, but the media is supposed to have has linguistic content"), nil
	}
	return report.Pass(), nil // no linguistical content
}
