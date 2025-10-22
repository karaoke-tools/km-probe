// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson/karamoe/songtype"
	"github.com/karaoke-tools/km-probe/internal/karajson/tag"
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe/baseprobe"
	"github.com/karaoke-tools/km-probe/internal/probes/report"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid/v5"
)

type NoOrigin struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewNoOrigin() probe.Probe {
	return &NoOrigin{
		baseprobe.New(
			"no-origin",
			"songtype is OP/ED/IN but origin tag is missing",
			cond.Any{
				cond.HasMoreTagsThan{
					TagType: tag.Origins,
					Number:  0,
					Msg:     "has origin",
				},
				cond.HasNoTagFrom{
					TagType: tag.Songtypes,
					Tags:    []uuid.UUID{songtype.OP, songtype.ED, songtype.IN},
					Msg:     "songtype is not OP/ED/IN",
				},
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p NoOrigin) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	return report.Fail(severity.Critical, "add the missing origin tag"), nil
}
