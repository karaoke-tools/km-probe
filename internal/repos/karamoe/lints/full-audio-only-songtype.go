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
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/songtype"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/version"

	"github.com/gofrs/uuid/v5"
)

type FullAudioOnlySongtype struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewFullAudioOnlySongtype() lint.Lint {
	return &FullAudioOnlySongtype{
		baselint.New(
			"full-audio-only-songtype",
			"audio only song are considered OP/ED/IN only when they have the same size than actual OP/ED/IN",
			cond.Any{
				cond.HasNoTagFrom{
					TagType: tag.Versions,
					Tags:    []uuid.UUID{version.Full},
					Msg:     "is not a full version",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.AUDIO},
					Msg:     "is not an audio only",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED, songtype.IN},
					Msg:     "is not an OP/ED/IN",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p FullAudioOnlySongtype) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "remove the OP/ED/IN tag"), nil
}
