// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package cliargs

import (
	"context"
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
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
	ErrUuidArgumentInvalid   = errors.New("Invalid value: KID argument is not matching the UUID format")
	ErrUnknownArgument       = errors.New("Unknown argument")
	ErrUnknownCommand        = errors.New("Unknown command")
)

// Usage error func
func UsageError(ctx context.Context, command *cli.Command, err error, isSubcommand bool) error {
	return err
}

// Errors when there are still unparsed arguments
func CheckUnknownArgs(ctx context.Context, command *cli.Command) (context.Context, error) {
	if command.Args().Present() {
		cli.ShowAppHelp(command)
		return ctx, ErrUnknownArgument
	}
	return ctx, nil
}

// Errors when the command does not exist
func CommandNotFound(ctx context.Context, command *cli.Command, s string) {
	cli.ShowAppHelp(command)
}

// Errors when flag does not exist
func InvalidFlagAccess(ctx context.Context, command *cli.Command, s string) {
	cli.ShowAppHelp(command)
}

// Validate argument type "WHEN"
func CheckWhen(ctx context.Context, command *cli.Command, v string) error {
	if !slices.Contains(when, v) {
		cli.ShowAppHelp(command)
		return ErrWhenArgumentInvalid
	}
	return nil
}

// Validate argument type "FORMAT"
func CheckFormat(ctx context.Context, command *cli.Command, v string) error {
	if !slices.Contains(format, v) {
		cli.ShowAppHelp(command)
		return ErrFormatArgumentInvalid
	}
	return nil
}

// Validate argument type "UUID"
func CheckUuids(ctx context.Context, command *cli.Command, v []string) error {
	for _, u := range v {
		if !re_uuid.Match([]byte(u)) {
			cli.ShowAppHelp(command)
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
