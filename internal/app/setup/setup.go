// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package setup

import (
	"github.com/urfave/cli/v2"
)

type Setup struct {
	// settings
	Hyperlink  bool
	Color      bool
	OutputJson bool
}

func FromCli(ctx *cli.Context) *Setup {
	s := &Setup{}

	// get value for hypelink
	switch ctx.String("hyperlink") {
	case "never":
		s.Hyperlink = false
	case "always":
		s.Hyperlink = true
	default: // auto
		// if on your platform/terminal you see strange symbols
		// please report this issue: https://github.com/louisroyer/km-probe/issues/new
		// TODO: disable if not a terminal
		s.Hyperlink = true
	}

	// get value for color, this enables the use of ansi codes
	switch ctx.String("color") {
	case "never":
		s.Color = false
	case "always":
		s.Color = true
	default: // auto
		// if on your platform/terminal you see strange symbols
		// please report this issue: https://github.com/louisroyer/km-probe/issues/new
		// TODO: disable if not a terminal
		s.Color = true
	}

	// get value for json
	switch ctx.String("output-format") {
	case "txt":
		s.OutputJson = false
	case "json":
		s.OutputJson = true
	default:
		// auto
		// TODO: disable if it is a pipe
		s.OutputJson = false
	}
	return s
}
