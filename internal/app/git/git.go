// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"errors"
	"path/filepath"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/app"
	"github.com/louisroyer/km-probe/internal/app/setup"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GitSetup struct {
	*setup.Setup
	Repositories []app.Repository
}

func FromCli(ctx *cli.Context) (*GitSetup, error) {
	s := &GitSetup{
		Setup:        setup.FromCli(ctx),
		Repositories: make([]app.Repository, 0),
	}

	kmConfig, err := app.LoadConf()
	if err != nil {
		return nil, err
	}
	for _, v := range kmConfig.System.Repositories {
		if repo := ctx.String("repository"); repo != "" && repo != v.Name {
			// we can only probe in the configured repository
			continue
		}
		baseDir, err := app.SearchKmDataDirPath(v.BaseDir)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"name":     v.Name,
				"base-dir": v.BaseDir,
			}).Error("Repository is configured with a base directory that doesn't exist")
			continue
		}

		mediaPath, err := app.SearchKmDataDirPath(v.Path.Medias[0]) // TODO: multi-track drifting
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"name":      v.Name,
				"media-dir": v.Path.Medias[0],
			}).Error("Repository is configured with a media directory that doesn't exist")
			continue
		}

		s.Repositories = append(s.Repositories, app.Repository{
			Name:      v.Name,
			BaseDir:   baseDir,
			MediaPath: mediaPath,
		})
	}
	return s, nil
}

func RunFromCli(ctx *cli.Context) error {
	if s, err := FromCli(ctx); err != nil {
		return err
	} else {
		return s.Run(ctx.Context)
	}
}

// Maximum number of karaokes processed simultaneously.
// It is not useful to increase this number
// because we are bound by the speed of the json encoder.
// Increasing the number of worker will consume more memory
// because we cannot recycle structures when they are still used;
// making the work of the garbage collector more difficult,
// which will slow down everything, and may make the interface irresponsible
// for enough time to be noticable (~1s).
const MAX_WORKERS = 0xFF

func (s *GitSetup) Run(ctx context.Context) error {
	printer := app.NewPrinter()

	// "modified" flag: it is set to true by any goroutine if a file is found
	// this is safe for concurrent use because when the value is updated, it is always to `true`
	modified := false

	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgRepos := sync.WaitGroup{}
	defer func() {
		wgRepos.Wait()
		if !modified {
			logrus.Info("All repositories are clean. No karaoke to probe.")
		}
	}()
	workers := make(chan struct{}, MAX_WORKERS) // maximum number of simultaneous workers
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wgRepos.Go(func() {
				// parse git status
				// for each modified kara
				kara, err := GitModifiedKaras(ctx, repo.BaseDir)
				if err != nil {
					if errors.Is(err, ErrUnresolvedMergeConflict) {
						logrus.WithError(err).Error("Merge conflict need to be resolved")
					}
					return
				}
				if len(kara) > 0 {
					modified = true
				}
				for _, u := range kara {
					p := filepath.Join(repo.BaseDir, "karaokes", u.String()+".kara.json")
					workers <- struct{}{}
					wg.Go(func() {
						app.RunOnFile(ctx, &repo, p, printer)
						<-workers
					})
				}
			})
		}
	}
	return nil
}
