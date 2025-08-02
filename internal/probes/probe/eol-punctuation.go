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
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type EolPunctuation struct {
	baseProbe
}

func NewEolPunctuation(karaData *karadata.KaraData) Probe {
	return &EolPunctuation{
		newBaseProbe("eol-punctuation", karaData),
	}
}

func (p *EolPunctuation) Run(ctx context.Context) (report.Report, error) {
	for _, line := range p.karaData.Lyrics.Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				l := line.Text.StripTags()
				if strings.HasSuffix(l, ".") || strings.HasSuffix(l, ",") {
					if strings.HasSuffix(l, "...") {
						continue
					}
					return report.Fail(severity.Critical, "remove useless punctuation (`.` or `,`) at end of line"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
