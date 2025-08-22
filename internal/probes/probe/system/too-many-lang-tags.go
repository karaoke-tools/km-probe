// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type TooManyLangTags struct {
	baseprobe.BaseProbe
}

func NewTooManyLangTags(karaData *karadata.KaraData) probe.Probe {
	return &TooManyLangTags{
		baseprobe.New("too-many-lang-tags",
			"if more than 2 langs tags, replace them with multilingual tag",
			cond.HasLessTagsThan{
				TagType: tag.Langs,
				Number:  3,
				Msg:     "has not more than 2 lang tags",
			},
			karaData),
	}
}

func (p *TooManyLangTags) Run(ctx context.Context) (report.Report, error) {
	return report.Fail(severity.Critical, "replace lang tags with \"multilingual\""), nil
}
