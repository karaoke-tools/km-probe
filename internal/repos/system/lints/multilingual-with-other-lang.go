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

type MultilingualWithOtherLang struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewMultilingualWithOtherLang() lint.Lint {
	return &MultilingualWithOtherLang{
		baselint.New("multilingual-with-other-lang",
			"if multilingual tag is applied, no other lang tag should be present",
			cond.Any{
				cond.HasLessTagsThan{
					TagType: tag.Langs,
					Number:  2,
					Msg:     "single lang tag",
				},
				cond.HasNoTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.MUL},
					Msg:     "has not multilingual tag",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p MultilingualWithOtherLang) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "check languages tags"), nil
}
