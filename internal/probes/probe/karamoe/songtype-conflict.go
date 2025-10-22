// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"context"
	"slices"

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

type SongtypeConflict struct {
	baseprobe.BaseProbe
	probe.WithDefault
}

func NewSongtypeConflict() probe.Probe {
	return &SongtypeConflict{
		baseprobe.New(
			"songtype-conflict",
			"detects incompatible songtypes",
			cond.HasLessTagsThan{
				TagType: tag.Songtypes,
				Number:  2,
				Msg:     "has a single songtype",
			},
		),
		baseprobe.EnabledByDefault{},
	}
}

func (p SongtypeConflict) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.OT}) {
		return report.Fail(severity.Critical, "songtype \"OT\" is forbidden"), nil
	}
	if KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.AUDIO}) && KaraData.KaraJson.HasAnyTagFrom(tag.Songtypes, []uuid.UUID{songtype.MV, songtype.AMV}) {
		return report.Fail(severity.Critical, "MV/AMV cannot be audio only"), nil
	}

	counter := 0
	for _, tag := range KaraData.KaraJson.Data.Tags.Songtypes {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// maybe in the future we will move AUDIO into "families", and forbid songs to be both a CS and something else (only allow audio only songs to be CS because"we are not an encyclopedia")
		// then we may simply force this field length to be equal to 1 (maybe directly in KM, by a rule in repo's manifest)
		if !slices.Contains([]uuid.UUID{songtype.AUDIO, songtype.CS}, tag) {
			counter++
			if counter > 1 {
				return report.Fail(severity.Critical, "incompatible songtypes"), nil
			}
		}
	}

	return report.Pass(), nil
}
