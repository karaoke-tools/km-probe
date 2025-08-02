// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/songtype"
	"github.com/louisroyer/km-probe/internal/karajson/warning"
	"github.com/louisroyer/km-probe/internal/probes/report"

	"github.com/gofrs/uuid"
)

type MediaWarningAudioOnly struct {
	baseProbe
}

func NewMediaWarningAudioOnly(karaData *karadata.KaraData) Probe {
	return &MediaWarningAudioOnly{
		newBaseProbe("media-warning-audio-only", karaData),
	}
}

// warnings that are related to the media
var mediaWarnings []uuid.UUID = []uuid.UUID{
	warning.R18Media,
	warning.Spoiler,
	warning.Epilepsy,
}

func (p *MediaWarningAudioOnly) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Skip(), nil
	}
	for _, w := range mediaWarnings {
		if slices.Contains(p.karaData.KaraJson.Data.Tags.Warnings, w) {
			return report.Fail(), nil
		}
	}
	return report.Pass(), nil
}
