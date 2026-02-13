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
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/version"

	"github.com/gofrs/uuid/v5"
)

type OffVocalWithoutParent struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewOffVocalWithoutParent() lint.Lint {
	return &OffVocalWithoutParent{
		baselint.New("off-vocal-without-parent",
			"off vocal but no parent",
			cond.HasNoTagFrom{
				TagType: tag.Versions,
				Tags:    []uuid.UUID{version.OffVocal},
				Msg:     "not an off vocal",
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p OffVocalWithoutParent) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Parents) == 0 {
		return report.Fail(severity.Critical, "add the right parent"), nil
	}

	return report.Pass(), nil
}
