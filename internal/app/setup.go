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
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/kmconfig"
	"github.com/louisroyer/km-probe/internal/probes"
)

type Setup struct {
	Repositories []Repository
}

func NewSetup() *Setup {
	setup := Setup{
		Repositories: make([]Repository, 0),
	}
	return &setup
}

func (s *Setup) Run(ctx context.Context) error {
	xdgDataHome, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		xdgDataHome = filepath.Join(usr.HomeDir, ".local/share")
	}
	xdgPath := filepath.Join(xdgDataHome, "karaokemugen-app/app/")
	kmConfig, err := kmconfig.ParseConf(filepath.Join(xdgPath, "config.yml"))
	if err != nil {
		return err
	}
	for _, v := range kmConfig.System.Repositories {
		baseDir := v.BaseDir
		if !filepath.IsAbs(baseDir) {
			baseDir = filepath.Join(xdgPath, baseDir)
		}
		mediaPath := v.Path.Medias[0] // TODO: why is it an array?
		if !filepath.IsAbs(mediaPath) {
			mediaPath = filepath.Join(xdgPath, mediaPath)
		}
		s.Repositories = append(s.Repositories, Repository{
			Name:      v.Name,
			BaseDir:   baseDir,
			MediaPath: mediaPath,
		})
	}

	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := filepath.WalkDir(filepath.Join(repo.BaseDir, "karaokes"), func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					if p == filepath.Join(repo.BaseDir, "karaokes") {
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
