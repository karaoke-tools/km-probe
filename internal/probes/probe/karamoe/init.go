// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"github.com/louisroyer/km-probe/internal/probes/probe"
)

// Available probes
func Probes() []probe.NewProbeFunc {
	return []probe.NewProbeFunc{
		NewAltVersionWithoutParent,
		NewCredits,
		NewDoubleConsonant,
		NewLiveDownload,
		NewLongTagOnShortMedia,
		NewMusicVideoCreditless,
		NewStyleSingleWhite,
		NewVowelMacron,
		NewWrongTsuSeparation,
	}
}
