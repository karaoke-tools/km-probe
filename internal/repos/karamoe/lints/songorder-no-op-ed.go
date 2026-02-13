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

	"github.com/gofrs/uuid/v5"
)

type SongorderNoOpEd struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewSongorderNoOpEd() lint.Lint {
	return &SongorderNoOpEd{
		baselint.New(
			"songorder-no-op-ed",
			"songorder is not compatible with this songtype",
			cond.Any{
				cond.HasNoSongorder{},
				cond.HasAnyTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED},
					Msg:     "songtype is an OP or ED",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p SongorderNoOpEd) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if b := KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.IN, songtype.PV}); b {
		return report.Fail(severity.Warning, "songorder in IS/PV may be justified, but is rare"), nil
	}
	// fanwork + MV/OT: yes this is a strange tag mix, but at least the "songorder" is probably a deliberate choice…
	// playlists might be better for that, but this is probably valid
	if b := KaraData.KaraJson.HasAnyTagFrom(tag.Origins, []uuid.UUID{origin.Fanworks}) && KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.MV, songtype.OT}); b {
		return report.Pass(), nil
	}
	// MV + OVA: probably a serie of music videos
	// playlists might be better for that, but this is probably valid
	if b := KaraData.KaraJson.HasAnyTagFrom(tag.Origins, []uuid.UUID{origin.OVA}) && KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.MV}); b {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "remove songorder, or add missing songtype"), nil
}
