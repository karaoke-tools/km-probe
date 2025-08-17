// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/songtype"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type AudioOnlyWithVideoContainer struct {
	baseprobe.BaseProbe
}

func NewAudioOnlyWithVideoContainer(karaData *karadata.KaraData) probe.Probe {
	return &AudioOnlyWithVideoContainer{
		baseprobe.New("audio-only-with-video-container",
			"audio only tag, but not an audio only media",
			karaData),
	}
}

var videoExtensions []string = []string{
	"mp4",
	"mkv",
	"webm",
}

func (p *AudioOnlyWithVideoContainer) Run(ctx context.Context) (report.Report, error) {
	filename := p.KaraData.KaraJson.Medias[0].Filename
	startExt := strings.LastIndexByte(filename, '.')
	extension := filename[startExt+1:]
	if !slices.Contains(videoExtensions, extension) {
		return report.Skip("is not an audio only"), nil
	}
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Fail(severity.Critical,
				"if this is a still image replace media with an audio container, "+
					"otherwise replace audio only tag with add appropriate family tag"),
			nil
	}
	return report.Pass(), nil
}
