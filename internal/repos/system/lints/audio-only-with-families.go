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
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/songtype"

	"github.com/gofrs/uuid/v5"
)

type AudioOnlyWithFamilies struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewAudioOnlyWithFamilies() lint.Lint {
	return &AudioOnlyWithFamilies{
		baselint.New("audio-only-with-families",
			"media content tag including both audio only tag and other tags at the same time",
			cond.HasNoTagFrom{
				TagType: tag.Songtypes,
				Tags:    []uuid.UUID{songtype.AudioOnly},
				Msg:     "not an audio only",
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p AudioOnlyWithFamilies) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Tags.Families) > 0 {
		return report.Fail(severity.Critical, "an audio only media cannot have a content type (family)"), nil
	}
	return report.Pass(), nil
}
