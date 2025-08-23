// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/version"
	"github.com/louisroyer/km-probe/internal/karajson/tag"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"

	"github.com/gofrs/uuid"
)

type OffVocalWithoutParent struct {
	baseprobe.BaseProbe
}

func NewOffVocalWithoutParent() probe.Probe {
	return &OffVocalWithoutParent{
		baseprobe.New("off-vocal-without-parent",
			"off vocal but no parent",
			cond.HasNoTagFrom{
				TagType: tag.Versions,
				Tags:    []uuid.UUID{version.OffVocal},
				Msg:     "not an off vocal",
			},
		),
	}
}

func (p OffVocalWithoutParent) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if len(KaraData.KaraJson.Data.Parents) == 0 {
		return report.Fail(severity.Critical, "add the right parent"), nil
	}

	return report.Pass(), nil
}
