// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/louisroyer/km-probe/internal/ass"
	"github.com/louisroyer/km-probe/internal/probe"

	"github.com/sirupsen/logrus"
)

type Setup struct {
	ConfigPath   string
	Repositories []Repository
}

func NewSetup() *Setup {
	return &Setup{
		ConfigPath: "/home/louis/local/share/karaokemugen-app/app/config.yml", // TODO: not hardcode
		Repositories: []Repository{
			// TODO: not hardcode
			Repository{
				Name:      "kara.moe",
				BaseDir:   "/home/louis/Documents/kara.moe",
				MediaPath: "repos/kara.moe/medias",
			},
		},
	}
}

func (s *Setup) Run(ctx context.Context) error {
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := filepath.Walk(path.Join(repo.BaseDir, "lyrics"), func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					fmt.Println(p)
					if p == path.Join(repo.BaseDir, "lyrics") {
						return nil
					}
					return filepath.SkipDir
				}
				f, err := os.OpenFile(p, os.O_RDONLY, 0)
				if err != nil {
					return err
				}
				defer f.Close()

				lyrics, err := ass.Parse(ctx, f)
				if err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"path": p,
					}).Error("Error parsing ass")
					return err
				}
				probe := probe.NewProbe(p, lyrics)
				if err := probe.Run(ctx); err != nil {
					return err
				}
				fmt.Println(probe.Report)
				return nil
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}
