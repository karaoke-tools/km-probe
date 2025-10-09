// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/karaoke-tools/km-probe/internal/app/cliargs"
	"github.com/karaoke-tools/km-probe/internal/app/git"
	"github.com/karaoke-tools/km-probe/internal/app/info"
	"github.com/karaoke-tools/km-probe/internal/app/karaokes"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	version := "Unknown version"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	app := &cli.Command{
		Name:                  "km-probe",
		Usage:                 "find common mistakes within your Karaoke Mugen repositories",
		EnableShellCompletion: true,
		Authors: []any{
			"Louis Royer",
		},
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "output-format",
				Usage:  strings.Join([]string{"configure output `FORMAT`; ", cliargs.DisplayFormat}, ""),
				Value:  "auto",
				Action: cliargs.CheckFormat,
			},
			&cli.StringFlag{
				Name:   "color",
				Usage:  strings.Join([]string{"colorize output `WHEN`; ", cliargs.DisplayWhen}, ""),
				Value:  "auto",
				Action: cliargs.CheckWhen,
			},
			&cli.StringFlag{
				Name:   "hyperlink",
				Usage:  strings.Join([]string{"create hyperlinks in output using OSC 8 escape sequence `WHEN`; only for non-json output; ", cliargs.DisplayWhen}, ""),
				Value:  "auto",
				Action: cliargs.CheckWhen,
			},
			&cli.StringSliceFlag{
				Name:     "repo",
				Usage:    "disable all repositories except this `REPOSITORY`",
				Required: false,
			},
		},
		OnUsageError:             cliargs.UsageError,
		InvalidFlagAccessHandler: cliargs.InvalidFlagAccess,
		CommandNotFound:          cliargs.CommandNotFound,
		Commands: []*cli.Command{
			{
				Name:                     "git",
				Usage:                    "Probes karaokes that has been modified locally and not yet committed to git",
				Before:                   cliargs.CheckUnknownArgs,
				Action:                   git.RunFromCommand,
				OnUsageError:             cliargs.UsageError,
				InvalidFlagAccessHandler: cliargs.InvalidFlagAccess,
			},
			{
				Name:                     "karaokes",
				Aliases:                  []string{"karaoke", "kara"},
				Usage:                    "Probes selected karaokes of all enabled repositories",
				Before:                   cliargs.CheckUnknownArgs,
				Action:                   karaokes.RunFromCommand,
				OnUsageError:             cliargs.UsageError,
				InvalidFlagAccessHandler: cliargs.InvalidFlagAccess,
				MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{{
					Flags: [][]cli.Flag{
						{
							&cli.StringSliceFlag{
								Name:     "kid",
								Usage:    "add karaokes with this `KID` (Karaoke UUID) to the selection",
								Required: false,
								Action:   cliargs.CheckUuids,
							},
						},
						{
							&cli.BoolFlag{
								Name:     "all",
								Usage:    "add all karaokes to the selection",
								Required: false,
							},
						},
					},
				}},
			},
			{
				Name:                     "info",
				Usage:                    "Shows a list of available probes",
				Before:                   cliargs.CheckUnknownArgs,
				Action:                   info.RunFromCommand,
				OnUsageError:             cliargs.UsageError,
				InvalidFlagAccessHandler: cliargs.InvalidFlagAccess,
			},
		},
	}
	if err := app.Run(ctx, os.Args); err != nil {
		logrus.WithError(err).Fatal("Fatal error while running the application")
	}
}
