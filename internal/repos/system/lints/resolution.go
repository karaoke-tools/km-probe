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

type Resolution struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewResolution() lint.Lint {
	return &Resolution{
		baselint.New("resolution",
			"resolution not set to 0×0",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p Resolution) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	if KaraData.Lyrics[0].ScriptInfo.PlayResX == 0 && KaraData.Lyrics[0].ScriptInfo.PlayResY == 0 {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "update resolution to be 0×0 (and check style size)"), nil
}
