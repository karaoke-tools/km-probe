// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/songtype"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type AudioOnlyWithFamilies struct {
	baseprobe.BaseProbe
}

func NewAudioOnlyWithFamilies(karaData *karadata.KaraData) probe.Probe {
	return &AudioOnlyWithFamilies{
		baseprobe.New("audio-only-with-families",
			"media content tag including both audio only tag and other tags at the same time",
			karaData),
	}
}

func (p *AudioOnlyWithFamilies) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Skip("not an audio only"), nil
	}
	if len(p.KaraData.KaraJson.Data.Tags.Families) > 0 {
		return report.Fail(severity.Critical, "an audio only media cannot have a content type (family)"), nil
	}
	return report.Pass(), nil
}
