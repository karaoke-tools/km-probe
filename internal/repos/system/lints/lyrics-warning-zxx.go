// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/warning"

	"github.com/gofrs/uuid/v5"
)

type LyricsWarningZXX struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewLyricsWarningZXX() lint.Lint {
	return &LyricsWarningZXX{
		baselint.New("lyrics-warning-zxx",
			"lyrics warning, but there is no linguistical content",
			cond.HasNoTagFrom{
				TagType: tag.Warnings,
				Tags:    []uuid.UUID{warning.R18Lyrics},
				Msg:     "no lyrics-warning tag",
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p LyricsWarningZXX) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Langs, language.ZXX) {
		return report.Fail(severity.Critical, "check if lyrics warning is relevant, and if the Langs field is set"), nil
	}
	return report.Pass(), nil
}
