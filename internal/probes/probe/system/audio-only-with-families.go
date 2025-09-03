// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/songtype"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type AudioOnlyWithFamilies struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewAudioOnlyWithFamilies() probe.Probe {
	return &AudioOnlyWithFamilies{
		baseprobe.New("audio-only-with-families",
			"media content tag including both audio only tag and other tags at the same time",
			cond.HasNoTagFrom{
				TagType: tag.Songtypes,
				Tags:    []uuid.UUID{songtype.AudioOnly},
				Msg:     "not an audio only",
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p AudioOnlyWithFamilies) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Tags.Families) > 0 {
		return report.Fail(severity.Critical, "an audio only media cannot have a content type (family)"), nil
	}
	return report.Pass(), nil
}
