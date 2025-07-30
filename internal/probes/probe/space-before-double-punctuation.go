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
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"

	"github.com/gofrs/uuid"
)

type SpaceBeforeDoublePunctuation struct {
	baseProbe
}

func NewSpaceBeforeDoublePunctuation(karaData *karadata.KaraData) Probe {
	return &SpaceBeforeDoublePunctuation{
		newBaseProbe("space-before-double-punctuation", karaData),
	}
}

func (p *SpaceBeforeDoublePunctuation) Run(ctx context.Context) (report.Report, error) {
	// we only check if language is full english, full japanese, or jpn+eng
	if len(p.karaData.KaraJson.Data.Tags.Langs) > 2 {
		return report.Skip(), nil
	} else if len(p.karaData.KaraJson.Data.Tags.Langs) == 2 {
		if !slices.Contains(p.karaData.KaraJson.Data.Tags.Langs, karajson.LangJPN) {
			return report.Skip(), nil
		}
		if !slices.Contains(p.karaData.KaraJson.Data.Tags.Langs, karajson.LangENG) {
			return report.Skip(), nil
		}
	} else if len(p.karaData.KaraJson.Data.Tags.Langs) == 1 {
		if !slices.Contains([]uuid.UUID{karajson.LangJPN, karajson.LangENG}, p.karaData.KaraJson.Data.Tags.Langs[0]) {
			return report.Skip(), nil
		}
	} else {
		// no lang, how is that possible?
		return report.Skip(), nil
	}

	for _, line := range p.karaData.Lyrics.Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && !((line.Type == lyrics.Comment) && (line.Effect != "karaoke")) {
				l := line.Text.StripTags()
				if strings.Contains(l, " ?") || strings.Contains(l, " !") ||
					strings.Contains(l, " ?") || strings.Contains(l, " !") { // non-breakable space
					return report.Fail(), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
