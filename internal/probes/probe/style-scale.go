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

type StyleScale struct {
	baseProbe
}

func NewStyleScale(karaData *karadata.KaraData) Probe {
	return &StyleScale{
		newBaseProbe("style-scale", karaData),
	}
}

func (p *StyleScale) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.Lyrics) == 0 {
		return report.Skip("no lyrics"), nil
	}

	// TODO: update this when multi-track drifting is released
	for _, line := range p.karaData.Lyrics[0].Styles {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return report.Abort(), err
				}
				if (s.ScaleX != "100") || (s.ScaleY != "100") {
					return report.Fail(severity.Critical, "check scale of styles"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
