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
	"github.com/louisroyer/km-probe/internal/karajson/system/warning"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type MediaWarningAudioOnly struct {
	baseprobe.BaseProbe
}

func NewMediaWarningAudioOnly(karaData *karadata.KaraData) probe.Probe {
	return &MediaWarningAudioOnly{
		baseprobe.New("media-warning-audio-only",
			"media warning but this is an audio only kara",
			karaData),
	}
}

// warnings that are related to the media
var mediaWarnings []uuid.UUID = []uuid.UUID{
	warning.R18Media,
	warning.Spoiler,
	warning.Epilepsy,
}

func (p *MediaWarningAudioOnly) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Skip("not an audio only"), nil
	}
	for _, w := range mediaWarnings {
		if slices.Contains(p.KaraData.KaraJson.Data.Tags.Warnings, w) {
			return report.Fail(severity.Critical, "check warning tags (maybe a R18-media should be changed to R18-lyrics, maybe a tag should be removed)"), nil
		}
	}
	return report.Pass(), nil
}
