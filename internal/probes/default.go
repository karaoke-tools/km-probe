// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

import (
	"github.com/karaoke-tools/km-probe/internal/probes/probe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/karamoe"
	"github.com/karaoke-tools/km-probe/internal/probes/probe/system"
)

var availableProbes = genAvailableProbes()
var enabledProbes = filterEnabled(availableProbes)

func genAvailableProbes() []probe.Probe {
	a := make([]probe.Probe, 0)
	a = append(a, system.Probes()...)
	a = append(a, karamoe.Probes()...)
	return a
}

func AvailableProbes() []probe.Probe {
	return availableProbes
}

func EnabledProbes() []probe.Probe {
	return enabledProbes
}

func filterEnabled(pl []probe.Probe) []probe.Probe {
	a := make([]probe.Probe, 0)
	for _, p := range pl {
		if p.Enabled() {
			a = append(a, p)
		}
	}
	return a
}
