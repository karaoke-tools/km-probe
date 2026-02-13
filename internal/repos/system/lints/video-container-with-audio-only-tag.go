// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"context"
	"slices"

	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
	"github.com/karaoke-tools/km-probe/internal/lints/report"
	"github.com/karaoke-tools/km-probe/internal/lints/report/severity"
	"github.com/karaoke-tools/km-probe/internal/lints/skip/cond"
	"github.com/karaoke-tools/km-probe/internal/repos/system/lints/baselint"
	"github.com/karaoke-tools/km-probe/internal/repos/system/tags/songtype"
)

type VideoContainerWithAudioOnlyTag struct {
	baselint.BaseLint
	lint.WithDefault
}

func NewVideoContainerWithAudioOnlyTag() lint.Lint {
	return &VideoContainerWithAudioOnlyTag{
		baselint.New("video-container-with-audio-only-tag",
			"video container, but audio only tag",
			cond.HasNotVideoExtension{},
		),
		baselint.EnabledByDefault{},
	}
}

func (p VideoContainerWithAudioOnlyTag) Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error) {
	if slices.Contains(KaraData.KaraJson.Data.Tags.Songtypes, songtype.AudioOnly) {
		return report.Fail(severity.Critical,
				"if this is a still image replace media with an audio container, "+
					"otherwise replace audio only tag with add appropriate family tag"),
			nil
	}
	return report.Pass(), nil
}
