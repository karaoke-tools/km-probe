// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson/misc"
	"github.com/louisroyer/km-probe/internal/karajson/origin"
	"github.com/louisroyer/km-probe/internal/karajson/songtype"
	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"

	"github.com/gofrs/uuid"
)

type Credits struct {
	baseProbe
}

var possibleCreditlessOrigin = []uuid.UUID{
	origin.Movie,
	origin.OriginalNetworkAnimation,
	origin.OriginalVideoAnimation,
	origin.TVSpecial,
	origin.TVSeries,
}

func isOriginCreditlessCompatible(o uuid.UUID) bool {
	return slices.Contains(possibleCreditlessOrigin, o)
}

func isOPED(st uuid.UUID) bool {
	return slices.Contains([]uuid.UUID{songtype.OP, songtype.ED}, st)
}

func NewCredits(karaData *karadata.KaraData) Probe {
	return &Credits{
		newBaseProbe("credits", karaData),
	}
}

func (p *Credits) Run(ctx context.Context) (report.Report, error) {
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Misc, misc.Creditless) {
		return report.Pass(), nil
	}
	if !slices.ContainsFunc(p.karaData.KaraJson.Data.Tags.Origins, isOriginCreditlessCompatible) {
		return report.Skip("origin not compatible with creditless"), nil
	}
	if !slices.ContainsFunc(p.karaData.KaraJson.Data.Tags.Songtypes, isOPED) {
		return report.Skip("songtype is not OP/ED"), nil
	}
	if slices.Contains(p.karaData.KaraJson.Data.Tags.Songtypes, songtype.AUDIO) {
		return report.Skip("audio only cannot be creditless"), nil
	}
	return report.Fail(severity.Warning, "if the media is already creditless, add the `Creditless`; if a creditless version exists (and is relevant!! see <https://kara.moe/playlist/quand-le-staff-fait-parti-du-generique> for counter-examples), update the media and add the tag"), nil
}
