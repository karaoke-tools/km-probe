// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"github.com/louisroyer/km-probe/internal/probes/analyser"
	"github.com/louisroyer/km-probe/internal/probes/probe"
)

var defaultProbes []probe.NewProbeFunc = []probe.NewProbeFunc{
	probe.NewAutomation,
	probe.NewDoubleConsonnant,
	probe.NewEolPunctuation,
	probe.NewSpaceBeforeDoublePunctuation,
	probe.NewLiveDownload,
	probe.NewResolution,
	probe.NewScaledBorderAndShadow,
	probe.NewSingleCollection,
	probe.NewStyleBlackBorder,
	probe.NewStyleScale,
	probe.NewStyleSingleWhite,
	probe.NewOffVocalWithoutParent,
	probe.NewKfShortSyllabes,
	probe.NewMediaWarningAudioOnly,
	probe.NewDoubleYearGroup,
	probe.NewWrongTsuSeparation,
	probe.NewLongTagOnShortMedia,
	probe.NewAudioOnlyWithVideoContainer,
}

var defaultAnalysers []analyser.NewAnalyserFunc = []analyser.NewAnalyserFunc{
	analyser.NewSuitableFirstContribution,
}
