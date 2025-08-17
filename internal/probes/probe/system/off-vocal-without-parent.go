// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/version"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type OffVocalWithoutParent struct {
	baseprobe.BaseProbe
}

func NewOffVocalWithoutParent(karaData *karadata.KaraData) probe.Probe {
	return &OffVocalWithoutParent{
		baseprobe.New("off-vocal-without-parent",
			"off vocal but no parent", karaData),
	}
}

func (p *OffVocalWithoutParent) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.KaraData.KaraJson.Data.Tags.Versions, version.OffVocal) {
		return report.Skip("not an off vocal"), nil
	}

	if len(p.KaraData.KaraJson.Data.Parents) == 0 {
		return report.Fail(severity.Critical, "add the right parent"), nil
	}

	return report.Pass(), nil
}
