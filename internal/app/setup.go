// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/kmconfig"
	"github.com/louisroyer/km-probe/internal/probes"

	"github.com/sirupsen/logrus"
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

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	ch := make(chan *probes.Aggregator)
	printCtx, cancelPrint := context.WithCancel(ctx)
	defer cancelPrint()
	generateCtx, cancelGenerate := context.WithCancel(printCtx)
	// not strictly required because we defer the CancelFunc of the parent, but `go vet` complains about it
	defer cancelGenerate()
	go func(ctx context.Context, cancel context.CancelFunc, ch <-chan *probes.Aggregator) {
		for {
			select {
			case <-ctx.Done():
				return
			case a, ok := <-ch:
				if !ok {
					cancel()
					return
				}
				if err := encoder.Encode(a); err != nil {
					logrus.WithError(err).Error("Error while encoding json")
				}
			}
		}
	}(printCtx, cancelPrint, ch)

	if err := func() error {
		defer close(ch)
		wg := sync.WaitGroup{}
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
					wg.Add(1)
					go func(ctx context.Context, p string) {
						defer wg.Done()
						karaJson, err := karajson.FromFile(p)
						if err != nil {
							select {
							case <-ctx.Done():
							default:
								logrus.WithError(err).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Could not parse karajson")
							}
							return
						}
						a, err := probes.FromKaraJson(ctx, repo.BaseDir, karaJson, nil)
						if errors.Is(err, karadata.ErrNoLyrics) {
							// skip
							return
						} else if err != nil {
							select {
							case <-ctx.Done():
							default:
								logrus.WithError(errors.Join(errors.New(d.Name()), err)).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Could not create probe aggregator")
							}
							return
						}
						if err := a.Run(ctx); err != nil {
							select {
							case <-ctx.Done():
							default:
								logrus.WithError(err).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Probe aggregator failure")
							}
							return
						} else {
							ch <- a
						}
					}(generateCtx, p)
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		wg.Wait()
		return nil
	}(); err != nil {
		return err
	}
	select {
	case <-generateCtx.Done():
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
}
