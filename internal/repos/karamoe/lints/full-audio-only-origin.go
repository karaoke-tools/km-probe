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
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/origin"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/songtype"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/version"

	"github.com/gofrs/uuid/v5"
)

type FullAudioOnlyOrigin struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewFullAudioOnlyOrigin() lint.Lint {
	return &FullAudioOnlyOrigin{
		baselint.New(
			"full-audio-only-origin",
			"audio only song cannot have origin when they have not the same size than actual OP/ED/IN",
			cond.Any{
				cond.HasEmptyTagtype{
					TagType: tag.Origins,
					Msg:     "has no origin",
				},
				cond.All{
					cond.HasLessTagsThan{
						TagType: tag.Origins,
						Number:  2,
						Msg:     "has a single origin tag",
					},
					cond.HasAnyTagFrom{
						TagType: tag.Origins,
						Tags:    []uuid.UUID{origin.Vtuber},
						Msg:     "vtuber tag",
					},
				},
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
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p FullAudioOnlyOrigin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "remove origin tag"), nil
}
