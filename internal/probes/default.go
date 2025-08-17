// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"github.com/louisroyer/km-probe/internal/probes/analyser"
	"github.com/louisroyer/km-probe/internal/probes/probe"
)

var defaultProbes = []probe.NewProbeFunc{
	probe.NewLiveDownload,
	probe.NewCredits,
	probe.NewGiantFont,
	probe.NewAutomation,                         // missing automation script
	probe.NewDoubleConsonnant,                   // double consonnant in same k-tag (jpn only)
	probe.NewEolPunctuation,                     // non-significant punctuation at end-of-lines
	probe.NewSpaceBeforeDoublePunctuation,       // space before double punctuation (jpn/eng only)
	probe.NewResolution,                         // resolusion not 0×0
	probe.NewScaledBorderAndShadow,              // scaled border and shadow not enabled
	probe.NewSingleCollection,                   // multiple collections
	probe.NewStyleBlackBorder,                   // non-black borders
	probe.NewStyleScale,                         // style with scaling parameter
	probe.NewStyleSingleWhite,                   // unfilled color is not white (only if single style)
	probe.NewOffVocalWithoutParent,              // off vocal but no parent
	probe.NewKfShortSyllabes,                    // kf on very short syllabes
	probe.NewMediaWarningAudioOnly,              // media warning but this is an audio only kara
	probe.NewLyricsWarningNoLinguisticalContent, //  lyrics warning, but there is no linguistical content
	probe.NewDoubleYearGroup,                    // multiple year groups tags
	probe.NewWrongTsuSeparation,                 // `t|su` separation is not correct
	probe.NewLongTagOnShortMedia,                // "long" tag should only be applied automatically
	probe.NewAudioOnlyWithVideoContainer,        // audio only tag, but not an audio only media
	probe.NewAudioOnlyWithFamilies,              // both audio only tag and media content tag at the same time
	probe.NewAltVersionWithoutParent,            // version tag, but there is no parent song
	probe.NewMusicVideoCreditless,               // MV with a creditless tag
	probe.NewVowelMacron,                        // ā, ē, ō, ī, ū in lyrics file
	probe.NewUnknownMediaContent,                // missing content tag
	probe.NewFullKf,                             // lyrics with a lot of kf
	probe.NewNoLyrics,                           // missing lyrics file
}

var defaultAnalysers = []analyser.NewAnalyserFunc{
	analyser.NewSuitableFirstContribution,
}
