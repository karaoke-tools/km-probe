// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type KTimed struct {
	baseprobe.BaseProbe
}

func NewKTimed() probe.Probe {
	return &KTimed{
		baseprobe.New("k-timed",
			"there is at least one k-tag in the lyrics file",
			cond.NoLyrics{},
		),
	}
}

func (p KTimed) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!(line.Type == lyrics.Comment && strings.HasPrefix(line.Effect, "template"))) {
				if len(line.Text.TagsSplit) > 1 {
					return report.Pass(), nil
				}
			}
		}
	}
	return report.Fail(severity.Critical, "karaoke must not simply be line-timed, they must be syllabe-timed"), nil
}
