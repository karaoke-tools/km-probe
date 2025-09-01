// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ass

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
)

type Ass struct {
	ScriptInfo              *ScriptInfo
	Styles                  []string
	Events                  []*lyrics.LyricsParser
	Extradata               []string
	AegisubGarbage          bool // has a Aegisub Project Garbage section
	Fonts                   bool // has Fonts section
	assUnknownSectionsCount int
}

func Parse(ctx context.Context, lrc io.Reader) (*Ass, error) {
	scanner := bufio.NewScanner(lrc)

	state := assInit
	ass := &Ass{
		ScriptInfo: &ScriptInfo{},
		Styles:     make([]string, 0),
		Events:     make([]*lyrics.LyricsParser, 0),
	}
	i := 0
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			i++
			line := scanner.Text()
			if i == 1 {
				// In some files, BOM is present multiple times for no reason
				BOM := string([]byte{0xEF, 0xBB, 0xBF})
				for strings.HasPrefix(line, BOM) {
					line = strings.TrimPrefix(line, BOM)
				}
			}
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				switch strings.TrimSuffix(strings.TrimPrefix(line, "["), "]") {
				case "Script Info":
					state = assScriptInfo
				case "Aegisub Project Garbage":
					state = assAegisubGarbage
					ass.AegisubGarbage = true
				case "V4+ Styles":
					state = assStyles
				case "Fonts":
					state = assFonts
					ass.Fonts = true
				case "Events":
					state = assEvents
				case "Aegisub Extradata":
					state = assAegisubExtradata
				default:
					state = assUnknownSection
					ass.assUnknownSectionsCount += 1
				}
				continue
			}
			if line == "" {
				// empty line
				continue
			}
			switch state {
			case assScriptInfo:
				if strings.HasPrefix(line, ";") {
					// comment line
					continue
				}
				lineSplit := strings.SplitN(line, ": ", 2)
				if len(lineSplit) != 2 {
					// unreadable
					continue
				}
				switch lineSplit[0] {
				case "PlayResX":
					res, err := strconv.ParseUint(lineSplit[1], 10, 32)
					if err != nil {
						return nil, ErrMalformedFile
					}
					ass.ScriptInfo.PlayResX = uint32(res)
				case "PlayResY":
					res, err := strconv.ParseUint(lineSplit[1], 10, 32)
					if err != nil {
						return nil, ErrMalformedFile
					}
					ass.ScriptInfo.PlayResY = uint32(res)
				case "ScaledBorderAndShadow":
					if err := ass.ScriptInfo.SetScaledBorderAndShadow(lineSplit[1]); err != nil {
						return nil, err
					}
				}
			case assStyles:
				ass.Styles = append(ass.Styles, line)
			case assEvents:
				lyr, err := lyrics.Parse(line)
				if err != nil {
					return nil, err
				}
				ass.Events = append(ass.Events, lyr)
			case assAegisubExtradata:
				ass.Extradata = append(ass.Extradata, line)
			default:
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ass, nil
}
