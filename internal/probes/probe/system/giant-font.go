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

type GiantFont struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewGiantFont() probe.Probe {
	return &GiantFont{
		baseprobe.New("giant-font",
			"fonts that have unusual big size",
			cond.NoLyrics{},
		),
		baseprobe.EnabledByDefault{},
	}
}

const (
	GIANT_FONT_SIZE_WARNING  = 27
	GIANT_FONT_SIZE_CRITICAL = 30
)

func (p GiantFont) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	warn := false
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") {
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if s.Fontsize >= GIANT_FONT_SIZE_CRITICAL {
					return report.Fail(severity.Critical, "found a style with a big fontsize: if resolution is already 0x0, consider reducing font size (it may be hard to identify big text as lyrics to actually sing)"), nil
				}
				if s.Fontsize >= GIANT_FONT_SIZE_WARNING {
					warn = true
				}
			}
		}
	}
	if warn {
		return report.Fail(severity.Warning, "found a style with a big fontsize: if resolution is already 0x0, consider reducing font size (it may be hard to identify big text as lyrics to actually sing) "), nil
	}
	return report.Pass(), nil
}
