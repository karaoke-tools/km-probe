// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/karajson"

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

func (s *Setup) RunSingle(ctx context.Context) (err error) {
	printer := NewPrinter()

	// found flag: it is set to true by any goroutine if a file is found
	// this is safe for concurrent use because when the value is updated, it is always to `true`
	found := false

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		if !found {
			// `err` is a named return value: this allow us to modify it inside the defer
			err = ErrKaraokeNotFound
			logrus.WithFields(logrus.Fields{
				"uuid": s.uuid,
			}).WithError(err).Error("Karaoke not found")
		}
	}()

	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Go(func() {
				fp := filepath.Join(repo.BaseDir, "karaokes", s.uuid.String()+".kara.json")
				err := runOnFile(ctx, &repo, fp, printer)
				if err == nil || !errors.Is(err, fs.ErrNotExist) {
					found = true
				}
			})
		}
	}
	return nil
}

// Run on all karaokes of all repositories
func (s *Setup) RunAll(ctx context.Context) error {
	printer := NewPrinter()
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgRepos := sync.WaitGroup{}
	defer wgRepos.Wait()
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wgRepos.Go(func() {
				repo.WalkKaraokes(ctx,
					func(ctx context.Context, r *Repository, p string) error {
						wg.Go(func() { runOnFile(ctx, r, p, printer) })
						return nil
					},
				)
			})
		}
	}
	return nil
}

// Parse the karaoke, run probes, and display result
// p is the filepath to the .kara.json file
func runOnFile(ctx context.Context, repo *Repository, p string, printer *Printer) error {
	karaJson, err := karajson.FromFile(ctx, p)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !errors.Is(err, fs.ErrNotExist) {
				logrus.WithError(err).WithFields(logrus.Fields{
					"repository": repo.Name,
					"filepath":   p,
				}).Error("Could not parse karajson")
			}
			return err
		}
	}
	karaData, err := karadata.FromKaraJson(ctx, repo.BaseDir, karaJson)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logrus.WithError(err).WithFields(logrus.Fields{
				"repository": repo.Name,
				"filepath":   p,
			}).Error("Could not create karadata")
			return err
		}
		return err
	}
	aggregator := printer.Aggregator()
	aggregator.Reset(repo.BaseDir, karaJson)
	if err := aggregator.Run(ctx, karaData); err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logrus.WithError(err).WithFields(logrus.Fields{
				"repository": repo.Name,
				"filepath":   p,
			}).Error("Probe aggregator failure")
			return err
		}
	}
	if err := printer.Encode(ctx, aggregator); err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logrus.WithError(err).WithFields(logrus.Fields{
				"repository": repo.Name,
				"filepath":   p,
			}).Error("Could not print aggregator result")
			return err
		}
	}
	return nil
}
