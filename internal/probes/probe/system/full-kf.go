// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"strings"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"
)

type FullKf struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewFullKf() probe.Probe {
	return &FullKf{
		baseprobe.New("full-kf",
			"lyrics with a lot of kf",
			cond.NoLyrics{},
		),
		baseprobe.EnabledByDefault{},
	}
}

const (
	// ratio of \kf over \k tags until it is critical
	FULL_KF_CRITICAL_RATIO_KF = 1
	FULL_KF_CRITICAL_RATIO_K  = 2 // 1/3 of tags are \kf

	FULL_KF_WARNING_RATIO_KF = 1
	FULL_KF_WARNING_RATIO_K  = 3 // 1/4 of tags are \kf
)

func (p FullKf) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	kf_count := 0
	k_count := 0
	// TODO: update this when multi-track drifting is released
	for _, line := range KaraData.Lyrics[0].Events {
		select {
		case <-ctx.Done():
			return report.Abort(), ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!(line.Type == lyrics.Comment && strings.HasPrefix(line.Effect, "template"))) {
				for _, syll := range line.Text.TagsSplit {
					select {
					case <-ctx.Done():
						return report.Abort(), ctx.Err()
					default:
						if strings.HasPrefix(syll, "{") {
							if strings.Contains(syll, "\\kf") {
								kf_count += 1
							} else {
								k_count += 1
							}
						}
					}
				}
			}
		}
	}
	if FULL_KF_CRITICAL_RATIO_KF*kf_count >= FULL_KF_CRITICAL_RATIO_K*k_count {
		return report.Fail(severity.Critical, "too many \\kf in the song"), nil
	}
	if FULL_KF_WARNING_RATIO_KF*kf_count >= FULL_KF_WARNING_RATIO_K*k_count {
		return report.Fail(severity.Warning, "too many \\kf in the song"), nil
	}
	return report.Pass(), nil
}
