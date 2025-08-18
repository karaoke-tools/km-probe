// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type StyleBlackBorder struct {
	baseprobe.BaseProbe
}

func NewStyleBlackBorder(karaData *karadata.KaraData) probe.Probe {
	return &StyleBlackBorder{
		baseprobe.New("style-black-border",
			"detects non-black border",
			cond.NoLyrics{},
			karaData),
	}
}

func (p *StyleBlackBorder) Run(ctx context.Context) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") { // we don't care about furigana styles for the now
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if s.OutlineColour != "&H00000000" {
					// border color must be black
					return report.Fail(severity.Warning, "outline must be black (cannot this probe can only check if this is pure black, nuances of black be be okay"), nil
				}
				break
			}
		}
	}
	return report.Pass(), nil
}
