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
		Action: func(ctx *cli.Context) error {
			if err := app.NewSetup().Run(ctx.Context); err != nil {
				logrus.WithError(err).Fatal("Error while running, exiting…")
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "parse-ass",
				Usage: "parse ASS file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Usage:    "File to parse",
						Required: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					if err := app.NewSetupParseAss(ctx.Path("file")).Run(ctx.Context); err != nil {
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
