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

type AegisubGarbage struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewAegisubGarbage() lint.Lint {
	return &AegisubGarbage{
		baselint.New("aegisub-garbage",
			"Aegisub Project Garbage section has not been removed",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p AegisubGarbage) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].AegisubGarbage {
		return report.Fail(severity.Info, "Aegisub Project Garbage section has not been removed; if you are integrating this song make sure to enable the \"cleanup lyrics\" function in Karaoke Mugen "), nil
	}
	return report.Pass(), nil
}
