// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/collection"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type DoubleConsonant struct {
	baseprobe.BaseProbe
}

func NewDoubleConsonant(karaData *karadata.KaraData) probe.Probe {
	return &DoubleConsonant{
		baseprobe.New(
			"double-consonant",
			"double consonant in same k-tag (JPN romaji only)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Collections,
					Tags:    []uuid.UUID{collection.Kana},
					Msg:     "kana karaoke",
				},
				cond.HasTagsNotFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.JPN},
					Msg:     "not a japanese only version",
				},
			},
			karaData),
	}
}

var doubleConsonants = []string{
	"kk",
	"gg",
	"ss",
	"zz",
	"tt",
	"dd",
	"nn",
	"bb",
	"pp",
	"mm",
	"rr",
}

func (p *DoubleConsonant) Run(ctx context.Context) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				save := ""
				for _, syll := range line.Text.TagsSplit {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						if !strings.HasPrefix(syll, "{") {
							if !strings.HasSuffix(save, " ") { // this is not a new word
								for _, double := range doubleConsonants {
									if strings.HasPrefix(syll, double) {
										return report.Fail(severity.Critical,
											"check for double consonants: "+
												"there is at least an uncorrectly splitted `"+
												syll+"`"), nil
									}
								}

							}
							save = syll
						}
					}

				}
			}
		}
	}
	return report.Pass(), nil
}
