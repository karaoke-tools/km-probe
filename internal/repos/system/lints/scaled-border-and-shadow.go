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

type ScaledBorderAndShadow struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewScaledBorderAndShadow() lint.Lint {
	return &ScaledBorderAndShadow{
		baselint.New("scaled-border-and-shadow",
			"scaled border and shadow not enabled",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p ScaledBorderAndShadow) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].ScriptInfo.ScaledBorderAndShadow {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "check the \"Scale border and shadow\" box"), nil
}
