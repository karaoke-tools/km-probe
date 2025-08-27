// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/louisroyer/km-probe/internal/app/cliargs"
	"github.com/louisroyer/km-probe/internal/app/info"
	"github.com/louisroyer/km-probe/internal/app/karaokes"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	app := &cli.App{
		Name:                 "KM-probe",
		Usage:                "Probe for karaoke quality",
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{Name: "Louis Royer"},
		},
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
		},
		Action: func(ctx *cli.Context) error {
			logrus.Info(ctx.Args())
			// XXX: workaround for https://github.com/urfave/cli/issues/1993
			// this disables the completion when the argument is `--`, which is better than bugged values
			if slices.Contains([]string{"--generate-bash-completion"}, ctx.Args().First()) {
				return nil
			}
			cli.ShowAppHelp(ctx)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "git",
				Usage: "[NOT IMPLEMENTED] Probes karaokes that has been modified locally and not yet committed to git",
			},
			{
				Name:    "karaokes",
				Aliases: []string{"karaoke", "kara"},
				Usage:   "Probes all karaokes of all enabled repositories",
				Before:  cliargs.CheckUnknownArgs,
				Action: func(ctx *cli.Context) error {
					if err := karaokes.RunFromCli(ctx); err != nil {
						logrus.WithError(err).Fatal("Error while running, exiting…")
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "uuid",
						Usage:    "select only karaokes with this `UUID`",
						Required: false,
						Action:   cliargs.CheckUuid},
					&cli.StringFlag{
						Name:     "repository",
						Usage:    "select only karaokes from this `REPOSITORY`",
						Required: false,
					},
				},
			},
			{
				Name:   "info",
				Usage:  "Shows a list of available probes",
				Before: cliargs.CheckUnknownArgs,
				Action: func(ctx *cli.Context) error {
					if err := info.RunFromCli(ctx); err != nil {
						logrus.WithError(err).Fatal("Error while running, exiting…")
					}
					return nil
				},
			},
		},
	}
	if err := app.RunContext(ctx, os.Args); err != nil {
		logrus.Fatal(err)
	}
}
