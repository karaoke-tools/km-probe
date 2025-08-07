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

type UnknownMediaContent struct {
	baseProbe
}

func NewUnknownMediaContent(karaData *karadata.KaraData) Probe {
	return &UnknownMediaContent{
		newBaseProbe("unknown-media-content", karaData),
	}
}

func (p *UnknownMediaContent) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.KaraJson.Data.Tags.Families) > 0 {
		return report.Pass(), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "indicate the media content type (animation, real, audio only)"), nil
}
