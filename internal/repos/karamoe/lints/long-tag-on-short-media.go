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
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/misc"
)

type LongTagOnShortMedia struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewLongTagOnShortMedia() lint.Lint {
	return &LongTagOnShortMedia{
		baselint.New("long-tag-on-short-media",
			"long tag added manually",
			cond.GreaterMediaDuration{Duration: 300},
		),
		baselint.EnabledByDefault{},
	}
}

func (p LongTagOnShortMedia) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Misc, misc.Long) {
		return report.Fail(severity.Critical, "remove long tag"), nil
	}
	return report.Pass(), nil
}
