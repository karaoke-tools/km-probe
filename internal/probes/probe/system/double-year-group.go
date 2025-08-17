// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/system/year"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type DoubleYearGroup struct {
	baseprobe.BaseProbe
}

func NewDoubleYearGroup(karaData *karadata.KaraData) probe.Probe {
	return &DoubleYearGroup{
		baseprobe.New("double-year-group",
			"double-year-group",
			karaData),
	}
}

// we cannot just checking the len of the group field,
// because on some repositories it is used for more than only years
var yearsGroup []uuid.UUID = []uuid.UUID{
	year.Y1950,
	year.Y1960,
	year.Y1970,
	year.Y1980,
	year.Y1990,
	year.Y2000,
	year.Y2010,
	year.Y2020,
}

func (p *DoubleYearGroup) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.KaraJson.Data.Tags.Groups) < 2 {
		return report.Pass(), nil
	}
	ok := false
	for _, group := range p.KaraData.KaraJson.Data.Tags.Groups {
		if found := slices.Contains(yearsGroup, group); found && ok {
			return report.Fail(severity.Critical, "remove all years group and save to apply (let years hooks do their job)"), nil
		} else if found {
			ok = true
		}
	}
	return report.Pass(), nil
}
