// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes/report"
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
	if !slices.Contains(p.karaData.KaraJson.Data.Tags.Versions, karajson.VersionOffVocal) {
		return report.Skip(), nil
	}

	if len(p.karaData.KaraJson.Data.Parents) == 0 {
		return report.Fail(), nil
	}

	return report.Pass(), nil
}
