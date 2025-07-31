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

	"github.com/gofrs/uuid"
)

type DoubleYearGroup struct {
	baseProbe
}

func NewDoubleYearGroup(karaData *karadata.KaraData) Probe {
	return &DoubleYearGroup{
		newBaseProbe("double-year-group", karaData),
	}
}

// warnings that are related to the media
var yearsGroup []uuid.UUID = []uuid.UUID{
	karajson.Years1950,
	karajson.Years1960,
	karajson.Years1970,
	karajson.Years1980,
	karajson.Years1990,
	karajson.Years2000,
	karajson.Years2010,
	karajson.Years2020,
}

func (p *DoubleYearGroup) Run(ctx context.Context) (report.Report, error) {
	if len(p.karaData.KaraJson.Data.Tags.Groups) < 2 {
		return report.Pass(), nil
	}
	ok := false
	for _, group := range p.karaData.KaraJson.Data.Tags.Groups {
		if found := slices.Contains(yearsGroup, group); found && ok {
			return report.Fail(), nil
		} else if found {
			ok = true
		}
	}
	return report.Pass(), nil
}
