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
	"path/filepath"
	"strings"
	"sync"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"
	"github.com/louisroyer/km-probe/internal/probes"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Setup struct {
	Repositories []Repository
	uuid         uuid.UUID
}

func NewSetup() *Setup {
	setup := Setup{
		Repositories: make([]Repository, 0),
	}
	return &setup
}

func (s *Setup) Run(ctx context.Context, id string) error {
	// parse uuid
	if id != "" {
		if u, err := uuid.FromString(id); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"uuid-argument": id,
			}).Error("Could not parse uuid")
			return err
		} else {
			s.uuid = u
		}
	}
	kmConfig, err := loadConf()
	if err != nil {
		return err
	}
	for _, v := range kmConfig.System.Repositories {
		baseDir, err := searchKmDataDirPath(v.BaseDir)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"name":     v.Name,
				"base-dir": v.BaseDir,
			}).Error("Repository is configured with a base directory that doesn't exist")
			continue
		}

		mediaPath, err := searchKmDataDirPath(v.Path.Medias[0]) // TODO: multi-track drifting
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"name":      v.Name,
				"media-dir": v.Path.Medias[0],
			}).Error("Repository is configured with a media directory that doesn't exist")
			continue
		}

		s.Repositories = append(s.Repositories, Repository{
			Name:      v.Name,
			BaseDir:   baseDir,
			MediaPath: mediaPath,
		})
	}
	if s.uuid.IsNil() {
		return s.RunAll(ctx)
	}
	return s.RunSingle(ctx)
}

func (s *Setup) RunSingle(ctx context.Context) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fp := filepath.Join(repo.BaseDir, "karaokes", s.uuid.String()+".kara.json")
			karaJson, err := karajson.FromFile(fp)
			if err == nil {
				a, err := probes.FromKaraJson(ctx, repo.BaseDir, karaJson)
				if errors.Is(err, karadata.ErrNoMedias) {
					return err
				} else if err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"filename": fp,
					}).Error("Could not create probe aggregator")
					return err
				}
				if err := a.Run(ctx); err != nil {
					logrus.WithError(err).WithFields(logrus.Fields{
						"filename": fp,
					}).Error("Probe aggregator failure")
					return err
				}
				if err := encoder.Encode(a); err != nil {
					logrus.WithError(err).Error("Error while encoding json")
				}
				return nil
			} else if errors.Is(err, fs.ErrNotExist) {
				continue
			} else {
				logrus.WithError(err).WithFields(logrus.Fields{
					"filename": fp,
				}).Error("Could not parse karajson")
				return err
			}

		}
	}
	err := ErrKaraokeNotFound
	logrus.WithFields(logrus.Fields{
		"uuid": s.uuid,
	}).WithError(err).Error("Karaoke not found")
	return err
}

func (s *Setup) RunAll(ctx context.Context) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	ch := make(chan *probes.Aggregator)

	printCtx, cancelPrint := context.WithCancel(ctx)
	defer cancelPrint()
	generateCtx, cancelGenerate := context.WithCancel(ctx)
	defer cancelGenerate()

	wg := sync.WaitGroup{}
	defer func() {
		// waiting in a gorouting to avoid waiting forever
		// if ctx is Done after entering defer
		ctxWait, cancel := context.WithCancel(ctx)
		go func(cancel context.CancelFunc) {
			wg.Wait()
			cancel()
		}(cancel)
		select {
		case <-ctx.Done():
		case <-ctxWait.Done():
		}
	}()

	go func(ctx context.Context, ch <-chan *probes.Aggregator) {
		for {
			select {
			case <-ctx.Done():
				return
			case a := <-ch:
				wg.Done()
				if err := encoder.Encode(a); err != nil {
					logrus.WithError(err).Error("Error while encoding json")
				}
			}
		}
	}(printCtx, ch)

	if err := func(ctx context.Context) error {
		for _, repo := range s.Repositories {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := filepath.WalkDir(filepath.Join(repo.BaseDir, "karaokes"), func(p string, d fs.DirEntry, err error) error {
					select {
					case <-ctx.Done():
						return filepath.SkipDir
					default:
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
					}
					go func(ctx context.Context, p string) {
						select {
						case <-ctx.Done():
							return
						default:
							wg.Add(1)
						}
						karaJson, err := karajson.FromFile(p)
						if err != nil {
							select {
							case <-ctx.Done():
								return
							default:
								wg.Done()
								logrus.WithError(err).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Could not parse karajson")
								return
							}
						}
						a, err := probes.FromKaraJson(ctx, repo.BaseDir, karaJson)
						if errors.Is(err, karadata.ErrNoMedias) {
							// skip
							wg.Done()
							return
						} else if err != nil {
							select {
							case <-ctx.Done():
								return
							default:
								wg.Done()
								logrus.WithError(errors.Join(errors.New(d.Name()), err)).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Could not create probe aggregator")
								return
							}
						}
						if err := a.Run(ctx); err != nil {
							select {
							case <-ctx.Done():
								return
							default:
								wg.Done()
								logrus.WithError(err).WithFields(logrus.Fields{
									"filename": p,
								}).Error("Probe aggregator failure")
								return
							}
						} else {
							select {
							case <-ctx.Done():
								return
							case ch <- a:
								return
							}
						}
					}(ctx, p)
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}(generateCtx); err != nil {
		return err
	}
	return nil
}
