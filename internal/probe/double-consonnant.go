// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karajson"
)

const checkNoDoubleConsonnantIssuesKey = "double-consonnant"

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

func (p *Probe) checkNoDoubleConsonnantIssues(ctx context.Context) error {
	// we only check if language is full jpn
	for _, tag := range p.KaraJson.Data.Tags.Langs {
		if tag != karajson.LangJPN {
			p.Report.Pass(checkNoDoubleConsonnantIssuesKey)
			return nil
		}
	}
	for _, line := range p.Lyrics.Events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				save := ""
				for _, syll := range line.Text.TagsSplit {
					if !strings.HasPrefix(syll, "{") {
						if !strings.HasSuffix(save, " ") { // this is not a new word
							for _, double := range doubleConsonnants {
								if strings.HasPrefix(syll, double) {
									p.Report.Fail(checkNoDoubleConsonnantIssuesKey)
									return nil
								}
							}

						}
						save = syll
					}

				}
			}
		}
	}
	p.Report.Pass(checkNoDoubleConsonnantIssuesKey)
	return nil
}
