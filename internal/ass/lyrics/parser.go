// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lyrics

import (
	"regexp"
	"strconv"
	"strings"
)

type Type int

const (
	Format Type = iota
	Comment
	Dialogue
)

type LyricsParser struct {
	Type    Type
	Layer   string
	Start   string
	End     string
	Style   string
	Name    string
	MarginL string
	MarginR string
	MarginV string
	Effect  string
	Text    Line
}

func Parse(lyrics string) (*LyricsParser, error) {
	r := strings.SplitN(lyrics, ": ", 2)
	if len(r) != 2 {
		return nil, ErrMalformedLine
	}
	var t Type
	switch r[0] {
	case "Format":
		t = Format
	case "Comment":
		t = Comment
	case "Dialogue":
		t = Dialogue
	default:
		return nil, ErrMalformedLine
	}
	r1 := strings.SplitN(r[1], ",", 10)
	if len(r1) != 10 {
		return nil, ErrMalformedLine
	}
	return &LyricsParser{
		Type:    t,
		Layer:   r1[0],
		Start:   r1[1],
		End:     r1[2],
		Style:   r1[3],
		Name:    r1[4],
		MarginL: r1[5],
		MarginR: r1[6],
		MarginV: r1[7],
		Effect:  r1[8],
		Text:    NewLine(r1[9]),
	}, nil
}

type Line struct {
	TagsSplit []string
}

type lineState int

const (
	lineStateText lineState = iota
	lineStateTag
)

func NewLine(text string) Line {
	var buff strings.Builder
	line := make([]string, 0)
	state := lineStateText
	for _, letter := range text {
		switch state {
		case lineStateText:
			if letter == rune('{') {
				state = lineStateTag
				if buff.Len() > 0 {
					line = append(line, buff.String())
					buff.Reset()
				}
			}
			buff.WriteRune(letter)

		case lineStateTag:
			buff.WriteRune(letter)
			if letter == rune('}') {
				state = lineStateText
				line = append(line, buff.String())
				buff.Reset()
			}
		}
	}
	if buff.Len() > 0 {
		line = append(line, buff.String())
	}
	return Line{
		TagsSplit: line,
	}
}

func (l Line) StripTags() string {
	r := make([]string, 0, len(l.TagsSplit)/2)
	for _, e := range l.TagsSplit {
		if !strings.HasPrefix(e, "{") {
			r = append(r, e)
		}
	}
	return strings.Join(r[:], "")
}

func (l Line) String() string {
	return strings.Join(l.TagsSplit[:], "")
}

var re_kf_len = regexp.MustCompile(`{\\kf(\d+)}\s*[^\s{]`)

func (l Line) KfLen() []int {
	r := make([]int, 0)
	if f := re_kf_len.FindAllStringSubmatch(l.String(), len(l.TagsSplit)/2); f != nil {
		for _, e := range f {
			if len(e) < 2 {
				continue
			}
			if i, err := strconv.Atoi(e[1]); err == nil {
				r = append(r, i)
			}
		}
	}
	return r
}

var re_style = regexp.MustCompile(`(?:\\)r(.*?)(?:\\|\})`)

func (l *LyricsParser) Styles() []string {
	r := make([]string, 1, 1)
	r[0] = l.Style
	for _, e := range l.Text.TagsSplit {
		if f := re_style.FindStringSubmatch(e); f != nil {
			r = append(r, f[1])
		}
	}
	return r

}
