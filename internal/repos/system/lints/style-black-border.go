// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/style"
	"github.com/karaoke-tools/km-probe/internal/ass/style/colour"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
)

type StyleBlackBorder struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewStyleBlackBorder() lint.Lint {
	return &StyleBlackBorder{
		baselint.New("style-black-border",
			"detects non-black border",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p StyleBlackBorder) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") { // we don't care about furigana styles for the now
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if s.OutlineColour != colour.Black {
					// border color must be black
					return report.Fail(severity.Warning, "outline must be black (this lint can only check if this is pure black, nuances of black might be okay"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
