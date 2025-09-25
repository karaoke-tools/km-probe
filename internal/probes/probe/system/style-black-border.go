// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/style"
	"github.com/karaoke-tools/km-probe/internal/ass/style/colour"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"
)

type StyleBlackBorder struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewStyleBlackBorder() probe.Probe {
	return &StyleBlackBorder{
		baseprobe.New("style-black-border",
			"detects non-black border",
			cond.NoLyrics{},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p StyleBlackBorder) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") { // we don't care about furigana styles for the now
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if s.OutlineColour != colour.Black {
					// border color must be black
					return report.Fail(severity.Warning, "outline must be black (this probe can only check if this is pure black, nuances of black might be okay"), nil
				}
				break
			}
		}
	}
	return report.Pass(), nil
}
