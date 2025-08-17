// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/collection"
	"github.com/louisroyer/km-probe/internal/karajson/language"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type DoubleConsonant struct {
	baseProbe
}

func NewDoubleConsonant(karaData *karadata.KaraData) Probe {
	return &DoubleConsonant{
		newBaseProbe("double-consonant", karaData),
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
	if len(p.karaData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	// we only check if language is full jpn romaji
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Collections, collection.Kana) {
		return report.Skip("kana version"), nil
	}
	if res, err := p.karaData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip("not a japanese only version"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.karaData.Lyrics[0].Events {
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
