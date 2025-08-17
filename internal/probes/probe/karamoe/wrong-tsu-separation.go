// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/karamoe/collection"
	"github.com/louisroyer/km-probe/internal/karajson/system/language"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type WrongTsuSeparation struct {
	baseprobe.BaseProbe
}

// In japanese, `つ` should be timed as a single syllabe: `tsu`.
// For example, `ひとつ`(`hitotsu`) should be timed as `hi|to|tsu` and not as `hi|tot|su`.
func NewWrongTsuSeparation(karaData *karadata.KaraData) probe.Probe {
	return &WrongTsuSeparation{
		baseprobe.New("wrong-tsu-separation",
			"`t|su` separation is not correct (JPN romaji only)",
			karaData),
	}
}

func (p *WrongTsuSeparation) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}
	// we only check if language is full jpn romaji
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Collections, collection.Kana) {
		// TODO: with multi-track drifting, maybe we will have a way to detect kana versions
		// with System tags? so we could put this as system probe…
		return report.Skip("kana karaoke"), nil
	}
	if res, err := p.KaraData.KaraJson.HasOnlyLanguagesFrom(ctx, []uuid.UUID{language.JPN}); err != nil {
		return report.Abort(), err
	} else if !res {
		return report.Skip("not a japanese karaoke"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				ok := false
				for _, syll := range line.Text.TagsSplit {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						if !strings.HasPrefix(syll, "{") {
							if strings.HasSuffix(syll, "t") {
								ok = true
							} else if ok && strings.HasPrefix(syll, "su") {
								return report.Fail(severity.Critical, "`tsu` must be timed as a single syllabe"), nil
							} else {
								ok = false
							}
						}
					}
				}
			}
		}
	}
	return report.Pass(), nil
}
