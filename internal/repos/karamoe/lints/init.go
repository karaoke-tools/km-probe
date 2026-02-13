// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"github.com/karaoke-tools/km-probe/internal/lints"
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
)

func init() {
	lints.Register([]lint.Lint{
		NewAltVersionWithoutParent(),
		NewAudioOnlyCreditless(),
		NewAutomationAppliedFurigana(),
		NewAutomationAppliedNoFurigana(),
		NewCreditless(),
		NewDoubleConsonant(),
		NewFullAudioOnlyOrigin(),
		NewFullAudioOnlySongtype(),
		NewLiveDownload(),
		NewLongTagOnShortMedia(),
		NewMusicVideoCreditless(),
		NewNoOrigin(),
		NewSongorderNoOpEd(),
		NewSongtypeConflict(),
		NewStyleSingleWhite(),
		NewVersionConflict(),
		NewVowelMacron(),
		NewWrongTsuSeparation(),
	})
}
