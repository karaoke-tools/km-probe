// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
)

type SingleCollection struct {
	baseProbe
}

func NewSingleCollection(karaData *karadata.KaraData) Probe {
	return &SingleCollection{
		newBaseProbe("single-collection", karaData),
	}
}

func (p *SingleCollection) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.KaraJson.Data.Tags.Collections) != 1 {
		return report.Fail(severity.Critical, "choose the right collection according to rules"), nil
	}
	return report.Pass(), nil
}
