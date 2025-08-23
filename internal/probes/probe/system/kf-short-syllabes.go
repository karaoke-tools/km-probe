// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"strconv"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type KfShortSyllabes struct {
	baseprobe.BaseProbe
}

func NewKfShortSyllabes() probe.Probe {
	return &KfShortSyllabes{
		baseprobe.New("kf-short-syllabes",
			"kf on very short syllabes",
			cond.NoLyrics{},
		),
	}
}

func (p KfShortSyllabes) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	warning := false
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
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
