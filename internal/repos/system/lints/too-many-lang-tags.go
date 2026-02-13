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
)

type TooManyLangTags struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewTooManyLangTags() lint.Lint {
	return &TooManyLangTags{
		baselint.New("too-many-lang-tags",
			"if more than 2 langs tags, replace them with multilingual tag",
			cond.HasLessTagsThan{
				TagType: tag.Langs,
				Number:  3,
				Msg:     "has not more than 2 lang tags",
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p TooManyLangTags) Run(ctx context.Context, karaData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "replace lang tags with \"multilingual\""), nil
}
