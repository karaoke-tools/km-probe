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

type UnicodeWeirdSpaces struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewUnicodeWeirdSpaces() lint.Lint {
	return &UnicodeWeirdSpaces{
		baselint.New("unicode-weird-spaces",
			"detect lyrics file with weird unicode spaces",
			cond.Any{
				cond.NoLyrics{},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p UnicodeWeirdSpaces) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && !((line.Type == lyrics.Comment) && (line.Effect != "karaoke")) {
				l := line.Text.StripTags()
				if strings.ContainsRune(l, '\u2005') { // FOUR-PER-EM SPACE
					return report.Fail(severity.Warning,
							"Found `Four-Per-Em Space` (U+2005): replace it with regular space "+
								"(they may not render correcly on all systems)",
						),
						nil
				}
			}
		}
	}
	return report.Pass(), nil
}
