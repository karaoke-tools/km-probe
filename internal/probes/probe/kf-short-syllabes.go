// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strconv"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type KfShortSyllabes struct {
	baseProbe
}

func NewKfShortSyllabes(karaData *karadata.KaraData) Probe {
	return &KfShortSyllabes{
		newBaseProbe("kf-short-syllabes", karaData),
	}
}

func (p *KfShortSyllabes) Run(ctx context.Context) (report.Report, error) {
	warning := false
	for _, line := range p.karaData.Lyrics.Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				for _, l := range line.Text.KfLen() {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						// I consider kf100 to be the optimal limit
						// kf90 can also be okay sometimes
						// but kf85 or lower is definitely bad
						if l < 85 && l != 0 { // kf0 is the same as k0
							return report.Fail(severity.Critical, "remove very short \\kf (found a `"+strconv.Itoa(l)+"`)"), nil
						} else if l < 90 && l != 0 {
							warning = true
						}
					}
				}
			}
		}
	}
	if warning {
		return report.Fail(severity.Warning, "check if \\kf under 90 are relevant"), nil
	}
	return report.Pass(), nil
}
