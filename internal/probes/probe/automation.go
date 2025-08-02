// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type Automation struct {
	baseProbe
}

func NewAutomation(karaData *karadata.KaraData) Probe {
	return &Automation{
		newBaseProbe("automation", karaData),
	}
}

func (p *Automation) Run(ctx context.Context) (report.Report, error) {
	for _, line := range p.karaData.Lyrics.Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if line.Type == lyrics.Comment {
				return report.Pass(), nil
			}
		}
	}
	return report.Fail(severity.Critical, "missing automation line in the lyrics file"), nil
}
