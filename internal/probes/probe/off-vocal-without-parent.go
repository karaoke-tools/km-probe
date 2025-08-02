// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/version"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type OffVocalWithoutParent struct {
	baseProbe
}

func NewOffVocalWithoutParent(karaData *karadata.KaraData) Probe {
	return &OffVocalWithoutParent{
		newBaseProbe("off-vocal-without-parent", karaData),
	}
}

func (p *OffVocalWithoutParent) Run(ctx context.Context) (report.Report, error) {
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Versions, version.OffVocal) {
		return report.Skip("not an off vocal"), nil
	}

	if len(p.karaData.KaraJson.Data.Parents) == 0 {
		return report.Fail(severity.Critical, "add the right parent"), nil
	}

	return report.Pass(), nil
}
