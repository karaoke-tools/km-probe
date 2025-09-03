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
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type VideoContainerWithAudioOnlyTag struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewVideoContainerWithAudioOnlyTag() probe.Probe {
	return &VideoContainerWithAudioOnlyTag{
		baseprobe.New("video-container-with-audio-only-tag",
			"video container, but audio only tag",
			cond.HasNotVideoExtension{},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p VideoContainerWithAudioOnlyTag) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Fail(severity.Critical,
				"if this is a still image replace media with an audio container, "+
					"otherwise replace audio only tag with add appropriate family tag"),
			nil
	}
	return report.Pass(), nil
}
