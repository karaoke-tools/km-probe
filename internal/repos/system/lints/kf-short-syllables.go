// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"strconv"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
)

type KfShortSyllables struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewKfShortSyllables() lint.Lint {
	return &KfShortSyllables{
		baselint.New("kf-short-syllables",
			"kf on very short syllables",
			cond.NoLyrics{},
		),
		baselint.EnabledByDefault{},
	}
}

const (
	// I consider kf100 to be the optimal limit kf90 can also be okay sometimes but kf75 or lower is definitely bad
	shortSyllableCriticalThreshold = 75
	shortSyllableWarningThreshold  = 90
)

func (p KfShortSyllables) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	hasKaraokeEffectLines := false
	warning := false
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (line.Type == lyrics.Comment) && (line.Effect == "karaoke") {
				hasKaraokeEffectLines = true
				for _, l := range line.Text.KfLen() {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						if l < shortSyllableCriticalThreshold && l != 0 { // kf0 is the same as k0
							return report.Fail(severity.Critical, "remove very short \\kf (found a `"+strconv.Itoa(l)+"`)"), nil
						} else if l < shortSyllableWarningThreshold && l != 0 {
							warning = true
						}
					}
				}
			}
		}
	}
	if !hasKaraokeEffectLines { // fallback, for old karaokes
		for _, line := range KaraData.Lyrics[0].Events {
			select {
			case <-ctx.Done():
				return report.Abort(), ctx.Err()
			default:
				if (line.Type != lyrics.Format) && (line.Type != lyrics.Comment) {
					hasKaraokeEffectLines = true
					for _, l := range line.Text.KfLen() {
						select {
						case <-ctx.Done():
							return report.Abort(), ctx.Err()
						default:
							if l < shortSyllableCriticalThreshold && l != 0 { // kf0 is the same as k0
								return report.Fail(severity.Critical, "remove very short \\kf (found a `"+strconv.Itoa(l)+"`)"), nil
							} else if l < shortSyllableWarningThreshold && l != 0 {
								warning = true
							}
						}
					}
				}
			}
		}
	}
	if warning {
		return report.Fail(severity.Warning, "check if \\kf under "+strconv.Itoa(shortSyllableWarningThreshold)+" are relevant"), nil
	}
	return report.Pass(), nil
}
