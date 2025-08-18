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

type UnknownMediaContent struct {
	baseprobe.BaseProbe
}

func NewUnknownMediaContent(KaraData *karadata.KaraData) probe.Probe {
	return &UnknownMediaContent{
		baseprobe.New("unknown-media-content",
			"missing content tag",
			cond.Never{},
			KaraData),
	}
}

func (p *UnknownMediaContent) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.KaraJson.Data.Tags.Families) > 0 {
		return report.Pass(), nil
	}
	if slices.Contains(p.KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Pass(), nil
	}
	return report.Fail(severity.Critical, "indicate the media content type (animation, real, audio only)"), nil
}
