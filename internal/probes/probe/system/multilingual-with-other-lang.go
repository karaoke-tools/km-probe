// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/system/language"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type MultilingualWithOtherLang struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewMultilingualWithOtherLang() probe.Probe {
	return &MultilingualWithOtherLang{
		baseprobe.New("multilingual-with-other-lang",
			"if multilingual tag is applied, no other lang tag should be present",
			cond.Any{
				cond.HasLessTagsThan{
					TagType: tag.Langs,
					Number:  2,
					Msg:     "single lang tag",
				},
				cond.HasNoTagFrom{
					TagType: tag.Langs,
					Tags:    []uuid.UUID{language.MUL},
					Msg:     "has not multilingual tag",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p MultilingualWithOtherLang) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "check languages tags"), nil
}
