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
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type AudioOnlyWithFamilies struct {
	baseProbe
}

func NewAudioOnlyWithFamilies(karaData *karadata.KaraData) Probe {
	return &AudioOnlyWithFamilies{
		newBaseProbe("audio-only-with-families", karaData),
	}
}

func (p *AudioOnlyWithFamilies) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Skip("not an audio only"), nil
	}
	if len(p.karaData.KaraJson.Data.Tags.Families) > 0 {
		return report.Fail(severity.Critical, "an audio only media cannot have a content type (family)"), nil
	}
	return report.Pass(), nil
}
