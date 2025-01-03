// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package style

import (
	"strings"
)

type Style struct {
	Name            string
	Fontname        string
	Fontsize        string
	PrimaryColour   string
	SecondaryColour string
	OutlineColour   string
	BackColour      string
	Bold            string
	Italic          string
	Underline       string
	StrikeOut       string
	ScaleX          string
	ScaleY          string
	Spacing         string
	Angle           string
	BorderStyle     string
	Outline         string
	Shadow          string
	Alignment       string
	MarginL         string
	MarginR         string
	MarginV         string
	Encoding        string
}

func Parse(style string) (*Style, error) {
	r := strings.SplitN(style, ",", 23)
	if len(r) != 23 {
		return nil, ErrMalformedLine
	}
	s := &Style{
		Name:            r[0],
		Fontname:        r[1],
		Fontsize:        r[2],
		PrimaryColour:   r[3],
		SecondaryColour: r[4],
		OutlineColour:   r[5],
		BackColour:      r[6],
		Bold:            r[7],
		Italic:          r[8],
		Underline:       r[9],
		StrikeOut:       r[10],
		ScaleX:          r[11],
		ScaleY:          r[12],
		Spacing:         r[13],
		Angle:           r[14],
		BorderStyle:     r[15],
		Outline:         r[16],
		Shadow:          r[17],
		Alignment:       r[18],
		MarginL:         r[19],
		MarginR:         r[20],
		MarginV:         r[21],
		Encoding:        r[22],
	}
	return s, nil
}
