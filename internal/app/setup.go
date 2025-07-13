// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes"
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
			err := filepath.WalkDir(path.Join(repo.BaseDir, "karaokes"), func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					if p == path.Join(repo.BaseDir, "karaokes") {
						return nil
					}
					return filepath.SkipDir
				}
				if !strings.HasSuffix(d.Name(), ".kara.json") {
					return nil
				}
				karaJson, err := karajson.FromFile(p)
				if err != nil {
					return err
				}
				a, err := probes.FromKaraJson(ctx, repo.BaseDir, karaJson, nil)
				if errors.Is(err, karadata.ErrNoLyrics) {
					// skip
					return nil
				} else if err != nil {
					return errors.Join(errors.New(d.Name()), err)
				}
				if err := a.Run(ctx); err != nil {
					return err
				}
				fmt.Println(a)
				return nil
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}
