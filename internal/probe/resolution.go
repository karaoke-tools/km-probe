// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strconv"
)

type Resolution struct {
	LyricsFile *Ass
}

func NewResolution(lyrics *Ass) *Resolution {
	return &Resolution{LyricsFile: lyrics}
}

func (p *Resolution) Run(ctx context.Context) (*Report, error) {
	report := NewReport()
	pass := false
	if p.LyricsFile.ScriptInfo.PlayResX == 0 && p.LyricsFile.ScriptInfo.PlayResY == 0 {
		pass = true
	}
	report.Content["pass"] = strconv.FormatBool(pass)
	return report, nil
}
