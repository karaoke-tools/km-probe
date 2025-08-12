// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type GiantFont struct {
	baseProbe
}

func NewGiantFont(karaData *karadata.KaraData) Probe {
	return &GiantFont{
		newBaseProbe("giant-font", karaData),
	}
}

const (
	GIANT_FONT_SIZE_WARNING  = 27
	GIANT_FONT_SIZE_CRITICAL = 30
)

func (p *GiantFont) Run(ctx context.Context) (report.Report, error) {
	warn := false
	for _, line := range p.karaData.Lyrics.Styles {
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
