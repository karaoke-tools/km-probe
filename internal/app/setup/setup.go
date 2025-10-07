// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package setup

import (
	"os"

	"github.com/moby/term"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Maximum number of karaokes processed simultaneously.
// It is not useful to increase this number
// because we are bound by the speed of the json encoder.
// Increasing the number of worker will consume more memory
// because we cannot recycle structures when they are still used;
// making the work of the garbage collector more difficult,
// which will slow down everything, and may make the interface irresponsible
// for enough time to be noticable (~1s).
const MAX_WORKERS = 0xFF

type Setup struct {
	// settings
	Hyperlink  bool
	Color      bool
	OutputJson bool

	// workers
	workers     chan struct{}
	withWorkers bool
}

// Use `StartWork` to wait for resources to become available before starting a goroutine.
// Don't forget to use `StopWork` when the work is done.
func (s *Setup) StartWork() {
	if s.withWorkers {
		s.workers <- struct{}{}
	}
}

// Use `StopWork` when a work is done. You must have called `StartWork` before.
func (s *Setup) StopWork() {
	if s.withWorkers {
		select {
		case <-s.workers:
		default:
			logrus.Warning("Setup.StopWork() called but work has not been started using Setup.StartWork() before")
		}
	}
}

func FromCli(ctx *cli.Context) *Setup {
	s := &Setup{
		withWorkers: true,
		workers:     make(chan struct{}, MAX_WORKERS), // maximum number of simultaneous workers
	}
	isTerminal := term.IsTerminal(os.Stdout.Fd())

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
		// by default, we display colors if this is a terminal
		// note: colors are currently not supported with json output
		s.Color = isTerminal
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
