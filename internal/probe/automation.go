// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strconv"
	"strings"
)

type Automation struct {
	LyricsFile *Ass
}

func NewAutomation(lyrics *Ass) *Automation {
	return &Automation{LyricsFile: lyrics}
}

func (p *Automation) Run(ctx context.Context) (*Report, error) {
	report := NewReport()
	pass := false

	for _, line := range p.LyricsFile.Events {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if strings.HasPrefix(line, "Comment: ") {
				pass = true
				break
			}
		}
	}
	report.Content["pass"] = strconv.FormatBool(pass)
	return report, nil
}
