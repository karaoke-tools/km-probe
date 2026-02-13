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
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/songtype"
)

type UnknownMediaContent struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewUnknownMediaContent() lint.Lint {
	return &UnknownMediaContent{
		baselint.New("unknown-media-content",
			"missing content tag",
			cond.Never{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p UnknownMediaContent) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Tags.Families) > 0 {
		return report.Pass(), nil
	}
	if slices.Contains(KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "indicate the media content type (animation, real, audio only)"), nil
}
