// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lints

import (
	"github.com/karaoke-tools/km-probe/internal/lints/lint"
)

var availableLints = make([]lint.Lint, 0)

func Register(lints []lint.Lint) {
	availableLints = append(availableLints, lints...)
}

func Available() []lint.Lint {
	return availableLints
}

func Enabled() []lint.Lint {
	return filterEnabled(availableLints)
}

func filterEnabled(lints []lint.Lint) []lint.Lint {
	a := make([]lint.Lint, 0)
	for _, l := range lints {
		if l.Enabled() {
			a = append(a, l)
		}
	}
	return a
}
