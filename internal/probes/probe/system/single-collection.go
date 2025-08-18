// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system/baseprobe"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/skip/cond"
)

type SingleCollection struct {
	baseprobe.BaseProbe
}

func NewSingleCollection(karaData *karadata.KaraData) probe.Probe {
	return &SingleCollection{
		baseprobe.New("single-collection",
			"multiple collections for a single karaoke",
			cond.Never{},
			karaData),
	}
}

func (p *SingleCollection) Run(ctx context.Context) (report.Report, error) {
	if len(p.KaraData.KaraJson.Data.Tags.Collections) != 1 {
		return report.Fail(severity.Critical, "choose the right collection according to rules"), nil
	}
	return report.Pass(), nil
}
