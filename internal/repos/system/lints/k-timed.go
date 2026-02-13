// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
)

type KTimed struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewKTimed() lint.Lint {
	return &KTimed{
		baselint.New("k-timed",
			"there is at least one k-tag in the lyrics file",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p KTimed) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!(line.Type == lyrics.Comment && strings.HasPrefix(line.Effect, "template"))) {
				if len(line.Text.TagsSplit) > 1 {
					return report.Pass(), nil
				}
			}
		}
	}
	return report.Fail(severity.Critical, "songs must not simply be line-timed, they must be syllable-timed"), nil
}
