// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package system

import (
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
)

// Available probes
func Probes() []probe.Probe {
	return []probe.Probe{
		NewAegisubGarbage(),
		NewAudioOnlyWithFamilies(),
		NewAutomation(),
		NewDoubleYearGroup(),
		NewEmbeddedFonts(),
		NewEolPunctuation(),
		NewFromDisplayType(),
		NewFullKf(),
		NewGiantFont(),
		NewKfShortSyllables(),
		NewKTimed(),
		NewLyricsWarningZXX(),
		NewMediaWarningAudioOnly(),
		NewMultilingualWithOtherLang(),
		NewNoLyrics(),
		NewOffVocalWithoutParent(),
		NewResolution(),
		NewScaledBorderAndShadow(),
		NewSingleCollection(),
		NewSpaceBeforeDoublePunctuation(),
		NewStyleBlackBorder(),
		NewStyleScale(),
		NewTooManyLangTags(),
		NewUnknownMediaContent(),
		NewVideoContainerWithAudioOnlyTag(),
	}
}
