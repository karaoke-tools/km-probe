// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karamoe

import (
	"github.com/louisroyer/km-probe/internal/probes/probe"
)

// Available probes
func Probes() []probe.Probe {
	return []probe.Probe{
		NewAltVersionWithoutParent(),
		NewAudioOnlyCreditless(),
		NewCreditless(),
		NewDoubleConsonant(),
		NewLiveDownload(),
		NewLongTagOnShortMedia(),
		NewMusicVideoCreditless(),
		NewNoOrigin(),
		NewStyleSingleWhite(),
		NewVowelMacron(),
		NewWrongTsuSeparation(),
	}
}
