// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
)

type Ass struct {
	ScriptInfo *ScriptInfo
	Styles     []string
	Events     []string
}

type ScriptInfo struct {
	PlayResX uint32
	PlayResY uint32
}

type assState int

const (
	assInit assState = iota
	assScriptInfo
	assAegisubGarbage
	assStyles
	assEvents
)

func ParseAss(ctx context.Context, lyrics io.Reader) (*Ass, error) {
	scanner := bufio.NewScanner(lyrics)
	state := assInit
	ass := &Ass{
		ScriptInfo: &ScriptInfo{},
		Styles:     make([]string, 0),
		Events:     make([]string, 0),
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			switch strings.TrimSuffix(strings.TrimPrefix(line, "["), "]") {
			case "Script Info":
				state = assScriptInfo
			case "Aegisub Project Garbage":
				state = assAegisubGarbage
			case "V4+ Styles":
				state = assStyles
			case "Events":
				state = assEvents
			}
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
					return nil, ErrMalformedAss
				}
				ass.ScriptInfo.PlayResX = uint32(res)

			case "PlayResY":
				res, err := strconv.ParseUint(lineSplit[1], 10, 32)
				if err != nil {
					return nil, ErrMalformedAss
				}
				ass.ScriptInfo.PlayResY = uint32(res)
			}
		case assStyles:
			ass.Styles = append(ass.Styles, line)
		case assEvents:
			ass.Events = append(ass.Events, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ass, nil
}
