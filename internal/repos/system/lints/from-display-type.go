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

type FromDisplayType struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewFromDisplayType() lint.Lint {
	return &FromDisplayType{
		baselint.New("from-display-type",
			"weird values for from-display-type",
			cond.Never{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p FromDisplayType) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	switch KaraData.KaraJson.Data.FromDisplayType {
	case "series":
		return report.Fail(severity.Warning,
			"from-display-type is manually set to series, but this is already the default"), nil
	case "singergroups":
		fallthrough
	case "singers":
		if len(KaraData.KaraJson.Data.Tags.Series) == 0 {
			return report.Fail(severity.Warning,
				"from-display-type is manually set to `"+KaraData.KaraJson.Data.FromDisplayType+"` but this is already the default"), nil
		}
		fallthrough
	case "songwriters":
		fallthrough
	case "franchises":
		fallthrough
	case "creators":
		fallthrough
	case "":
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "from-display-type is set with a weird value: `"+KaraData.KaraJson.Data.FromDisplayType+"`"), nil
}
