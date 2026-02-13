// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
)

type EmbeddedFonts struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewEmbeddedFonts() lint.Lint {
	return &EmbeddedFonts{
		baselint.New("embedded-fonts",
			"lyrics file embeds fonts",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p EmbeddedFonts) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].Fonts {
		return report.Fail(severity.Critical, "lyrics file embeds fonts; consider using standard fonts instead because fonts embedding creates big lyrics file"), nil
	}
	return report.Pass(), nil
}
