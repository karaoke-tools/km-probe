// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type AudioOnlyWithVideoContainer struct {
	baseProbe
}

func NewAudioOnlyWithVideoContainer(karaData *karadata.KaraData) Probe {
	return &AudioOnlyWithVideoContainer{
		newBaseProbe("audio-only-with-video-container", karaData),
	}
}

var videoExtensions []string = []string{
	"mp4",
	"mkv",
	"webm",
}

func (p *AudioOnlyWithVideoContainer) Run(ctx context.Context) (report.Report, error) {
	filename := p.karaData.KaraJson.Medias[0].Filename
	startExt := strings.LastIndexByte(filename, '.')
	extension := filename[startExt+1:]
	if !slices.Contains(videoExtensions, extension) {
		return report.Skip(), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, karajson.TypeAudioOnly) {
		return report.Fail(), nil
	}
	return report.Pass(), nil
}
