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

	"github.com/gofrs/uuid"
)

type DoubleConsonnant struct {
	baseProbe
}

func NewDoubleConsonnant(karaData *karadata.KaraData) Probe {
	return &DoubleConsonnant{
		newBaseProbe("double-consonnant", karaData),
	}
}

var doubleConsonnants = []string{
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

func (p *DoubleConsonnant) Run(ctx context.Context) (report.Report, error) {
	// we only check if language is full jpn romaji
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Collections, collection.Kana) {
		return report.Skip(), nil
	}
	if res, err := p.karaData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip(), nil
	}

	for _, line := range p.karaData.Lyrics.Events {
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
								for _, double := range doubleConsonnants {
									if strings.HasPrefix(syll, double) {
										return report.Fail(), nil
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
