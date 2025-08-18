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

type EolPunctuation struct {
	baseprobe.BaseProbe
}

func NewEolPunctuation(karaData *karadata.KaraData) probe.Probe {
	return &EolPunctuation{
		baseprobe.New("eol-punctuation",
			"non-significant punctuation at end-of-lines",
			cond.NoLyrics{},
			karaData),
	}
}

func (p *EolPunctuation) Run(ctx context.Context) (report.Report, error) {
	// TODO: update this when multi-track drifting is released
	for _, line := range p.KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				l := line.Text.StripTags()
				if strings.HasSuffix(l, ".") || strings.HasSuffix(l, ",") {
					if strings.HasSuffix(l, "...") {
						continue
					}
					return report.Fail(severity.Critical, "remove useless punctuation (`.` or `,`) at end of line"), nil
				}
			}
		}
	}
	return report.Pass(), nil
}
