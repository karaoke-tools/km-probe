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

type SingleCollection struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewSingleCollection() lint.Lint {
	return &SingleCollection{
		baselint.New("single-collection",
			"multiple collections set",
			cond.Never{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p SingleCollection) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Tags.Collections) != 1 {
		return report.Fail(severity.Critical, "choose the right collection according to rules"), nil
	}
	return report.Pass(), nil
}
