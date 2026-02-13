// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/ass/style"
	"github.com/karaoke-tools/km-probe/internal/ass/style/colour"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/misc"
	"github.com/karaoke-tools/km-probe/internal/repos/karamoe/tags/origin"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/language"

	"github.com/gofrs/uuid/v5"
)

type StyleSingleWhite struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewStyleSingleWhite() lint.Lint {
	return &StyleSingleWhite{
		baselint.New("style-single-white",
			"unfilled color is not white (only if single style)",
			cond.Any{
				cond.NoLyrics{},
				cond.HasAnyTagFrom{
					TagType: tag.Misc,
					Tags:    []uuid.UUID{misc.GroupSinging},
					Msg:     "group singing song", // we can use one color by voice
				},
				cond.HasMoreTagsThan{
					TagType: tag.Langs,
					Number:  1,
					Msg:     "is multilingual song", // we can use one color by language
				},
				cond.HasAnyTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.MUL},
					Msg:     "is multilingual song", // we can use one color by language
				},
				cond.HasAnyTagFrom{
					TagType: tag.Origins,
					Tags:    []uuid.UUID{origin.Musical},
					// musicals often have dialogues,
					// and it's okay to have multiple colors even if this is not a group singing song
					Msg: "is from a musical",
				},
			},
		),
		baselint.EnabledByDefault{},
	}
}

func (p StyleSingleWhite) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	nonWhiteChoirStyleCnt := 0 // detected choir: update secondary color (if non-group kara)
	whiteChoirStyleCnt := 0    // detected choir: white secondary color
	nonWhiteUnknownStyleCnt := 0
	whiteUnknownStyleCnt := 0
	unused := 0

	// TODO: update this when multi-track drifting is released
	styles := make([]string, 0, len(KaraData.Lyrics[0].Styles)-1)

	// list of used styles
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if line.Type != lyrics.Dialogue {
				continue
			}
			for _, style := range line.Styles() {
				select {
				case <-ctx.Done():
					return report.Abort(), ctx.Err()
				default:
					if !slices.Contains(styles, style) {
						styles = append(styles, style)
					}
				}
			}
		}

	}

	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if !strings.HasPrefix(line, "Style: ") {
				// ignore format line
				continue
			}
			s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
			if err != nil {
				return report.Abort(), err
			}
			if !slices.Contains(styles, s.Name) {
				// unused style
				unused++
				continue
			}
			l_name := strings.ToLower(s.Name)
			if strings.Contains(l_name, "-furigana") {
				continue
			}
			choir := strings.Contains(l_name, "choir") ||
				strings.Contains(l_name, "spoken") ||
				strings.Contains(l_name, "dialogue") ||
				strings.Contains(l_name, "rubyscript") // when mixing rubscript lines with normal lines (2 template scripts)
			if s.SecondaryColour == colour.White {
				if choir {
					whiteChoirStyleCnt++
				} else {
					whiteUnknownStyleCnt++
				}
			} else {
				if choir {
					nonWhiteChoirStyleCnt++
				} else {
					nonWhiteUnknownStyleCnt++
				}
			}
		}
	}
	if nonWhiteUnknownStyleCnt == 1 {
		return report.Fail(severity.Critical, "update style: secondary color must be white"), nil
	}
	if nonWhiteUnknownStyleCnt > 1 {
		return report.Fail(severity.Critical, "if this is a group song, add the \"group-singing\" tag; if this is a multi-lingual song (with a color per language), add the missing lang tag; else make all secondary colors white"), nil
	}
	if nonWhiteChoirStyleCnt == 1 {
		return report.Fail(severity.Warning, "consider updating the secondary color of the choir style to white"), nil
	}
	if whiteUnknownStyleCnt > 1 {
		return report.Fail(severity.Warning, "found multiple styles with white as secondary color: should it be converted to group singing song?"), nil
	}
	// TODO: update this when multi-track drifting is released
	if unused > 0 {
		return report.Fail(severity.Info, "found some styles not used by Dialogue lines: if you are integrating this song make sure to enable the \"cleanup lyrics\" function in Karaoke Mugen (or maybe a style is used only by a Comment line?)"), nil
	}

	return report.Pass(), nil
}
