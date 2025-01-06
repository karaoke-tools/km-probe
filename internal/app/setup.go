// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probe"
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
			err := filepath.Walk(path.Join(repo.BaseDir, "karaokes"), func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					if p == path.Join(repo.BaseDir, "karaokes") {
						return nil
					}
					return filepath.SkipDir
				}
				if !strings.HasSuffix(info.Name(), ".kara.json") {
					return nil
				}
				karaJson, err := karajson.FromFile(p)
				if err != nil {
					return err
				}
				prb, err := probe.FromKaraJson(ctx, repo.BaseDir, karaJson)
				if errors.Is(err, probe.ErrNoLyrics) {
					// skip
					return nil
				} else if err != nil {
					return errors.Join(errors.New(info.Name()), err)
				}
				if err := prb.Run(ctx); err != nil {
					return err
				}
				fmt.Println(prb.Report)
				return nil
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}
