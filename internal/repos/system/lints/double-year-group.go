// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/year"

	"github.com/gofrs/uuid/v5"
)

type DoubleYearGroup struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewDoubleYearGroup() lint.Lint {
	return &DoubleYearGroup{
		baselint.New("double-year-group",
			"double-year-group",
			cond.Never{},
		),
		baselint.EnabledByDefault{},
	}
}

// we cannot just checking the len of the group field,
// because on some repositories it is used for more than only years
var yearsGroup []uuid.UUID = []uuid.UUID{
	year.Y1950,
	year.Y1960,
	year.Y1970,
	year.Y1980,
	year.Y1990,
	year.Y2000,
	year.Y2010,
	year.Y2020,
}

func (p DoubleYearGroup) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Tags.Groups) < 2 {
		return report.Pass(), nil
	}
	ok := false
	for _, group := range KaraData.KaraJson.Data.Tags.Groups {
		if found := slices.Contains(yearsGroup, group); found && ok {
			return report.Fail(severity.Critical, "remove all years group and save to apply (let years hooks do their job)"), nil
		} else if found {
			ok = true
		}
	}
	return report.Pass(), nil
}
