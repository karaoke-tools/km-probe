// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"github.com/louisroyer/km-probe/internal/probes/analyser"
	"github.com/louisroyer/km-probe/internal/probes/probe"
	"github.com/louisroyer/km-probe/internal/probes/probe/karamoe"
	"github.com/louisroyer/km-probe/internal/probes/probe/system"
)

var availableProbes = genAvailableProbes()

func genAvailableProbes() []probe.Probe {
	a := make([]probe.Probe, 0)
	a = append(a, system.Probes()...)
	a = append(a, karamoe.Probes()...)
	return a
}

func AvailableProbes() []probe.Probe {
	return availableProbes
}

var defaultAnalysers = []analyser.NewAnalyserFunc{
	analyser.NewSuitableFirstContribution,
}
