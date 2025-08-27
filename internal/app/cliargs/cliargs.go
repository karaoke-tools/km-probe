// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cliargs

import (
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
)

// Arguments types
var (
	when   = []string{"auto", "always", "never"}
	format = []string{"auto", "txt", "json"}
)

var (
	DisplayFormat = options("FORMAT", format)
	DisplayWhen   = options("WHEN", when)
)

var re_uuid = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// Argument parsing errors
var (
	ErrWhenArgumentInvalid   = errors.New(strings.Join([]string{"Invalid value: ", DisplayWhen}, ""))
	ErrFormatArgumentInvalid = errors.New(strings.Join([]string{"Invalid value: ", DisplayFormat}, ""))
	ErrUuidArgumentInvalid   = errors.New("Invalid value: UUID argument is not matching the UUID format")
	ErrUnknownArgument       = errors.New("Unknown argument")
)

// Errors when there are still unparsed arguments
func CheckUnknownArgs(ctx *cli.Context) error {
	if ctx.Args().Present() {
		// XXX: workaround for https://github.com/urfave/cli/issues/1993
		// this disables the completion when the argument is `--`, which is better than bugged values
		if !slices.Contains([]string{"--generate-bash-completion"}, ctx.Args().First()) {
			cli.ShowAppHelp(ctx)
		}
		return ErrUnknownArgument
	}
	return nil
}

// Validate argument type "WHEN"
func CheckWhen(ctx *cli.Context, v string) error {
	if !slices.Contains(when, v) {
		cli.ShowAppHelp(ctx)
		return ErrWhenArgumentInvalid
	}
	return nil
}

// Validate argument type "FORMAT"
func CheckFormat(ctx *cli.Context, v string) error {
	if !slices.Contains(format, v) {
		cli.ShowAppHelp(ctx)
		return ErrFormatArgumentInvalid
	}
	return nil
}

// Validate argument type "UUID"
func CheckUuids(ctx *cli.Context, v []string) error {
	for _, u := range v {
		if !re_uuid.Match([]byte(u)) {
			return ErrUuidArgumentInvalid
		}
	}
	return nil
}

// Get the options for this argument type
func options(name string, opts []string) string {
	tmp := make([]string, len(opts))
	for i, v := range opts {
		tmp[i] = strings.Join([]string{"\"", v, "\""}, "")
	}
	return strings.Join([]string{name, "=[", strings.Join(tmp, " | "), "]"}, "")
}
