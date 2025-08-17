// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

import (
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
)

type BaseProbe struct {
	pkg  string
	name string
	desc string
	// The following should only be used by the probe itself
	// to implement the `Run()` method
	KaraData *karadata.KaraData
}

func New(pkg string, name string, desc string, karaData *karadata.KaraData) BaseProbe {
	return BaseProbe{
		name:     strings.Join([]string{pkg, name}, "."),
		desc:     desc,
		KaraData: karaData,
	}
}

func (p *BaseProbe) Name() string {
	return p.name
}

func (p *BaseProbe) Description() string {
	return p.desc
}
