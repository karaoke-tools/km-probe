// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package setup

import (
	"os"

	"github.com/mattn/go-isatty"
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
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// get value for json
	switch ctx.String("output-format") {
	case "txt":
		s.OutputJson = false
	case "json":
		s.OutputJson = true
	default:
		s.OutputJson = !isTerminal
	}

	// get value for color, this enables the use of ansi codes
	switch ctx.String("color") {
	case "never":
		s.Color = false
	case "always":
		s.Color = true
	default:
		// by default, we display colors if this is not a json output
		// note: colors are currently not supported with json output
		s.Color = !s.OutputJson
		if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
			s.Color = false
		}
	}

	// get value for hypelink
	switch ctx.String("hyperlink") {
	case "never":
		s.Hyperlink = false
	case "always":
		s.Hyperlink = true
	default:
		// if there is color, there is no reason to not display hyperlinks by default
		s.Hyperlink = s.Color
	}

	return s
}
