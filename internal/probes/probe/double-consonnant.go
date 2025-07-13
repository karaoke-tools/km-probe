// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"
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
	// we only check if language is full jpn
	for _, tag := range p.karaData.KaraJson.Data.Tags.Langs {
		if tag != karajson.LangJPN {
			return report.Pass(), nil
		}
	}
	for _, line := range p.karaData.Lyrics.Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				save := ""
				for _, syll := range line.Text.TagsSplit {
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
	return report.Pass(), nil
}
