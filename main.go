// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/louisroyer/km-probe/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func checkUnknownArgs(ctx *cli.Context) error {
	if ctx.Args().Len() > 0 {
		cli.ShowAppHelp(ctx)
		logrus.WithFields(logrus.Fields{
			"unknown-args": ctx.Args(),
		}).Error("Unknown arguments")
	}
	return nil
}

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

		Commands: []*cli.Command{
			{
				Name:   "list-probes",
				Usage:  "info about available probes",
				Before: checkUnknownArgs,
				Action: func(ctx *cli.Context) error {
					if err := app.NewListProbes().Run(ctx.Context); err != nil {
						logrus.WithError(err).Fatal("Error while running, exiting…")
					}
					return nil
				},
			},
			{
				Name:   "run",
				Usage:  "info about available probes",
				Before: checkUnknownArgs,
				Action: func(ctx *cli.Context) error {
					if err := app.NewSetup().Run(ctx.Context, ctx.String("uuid")); err != nil {
						logrus.WithError(err).Fatal("Error while running, exiting…")
					}
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "uuid",
						Usage:    "UUID of a karaoke to probe",
						Required: false,
					},
				},
			},
		},
	}
	if err := app.RunContext(ctx, os.Args); err != nil {
		logrus.Fatal(err)
	}
}
