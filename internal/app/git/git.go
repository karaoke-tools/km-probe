// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"sync/atomic"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/app"
	"github.com/louisroyer/km-probe/internal/app/printer"
	"github.com/louisroyer/km-probe/internal/app/setup"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type GitSetup struct {
	*setup.Setup
	Repositories []app.Repository
	BaseUri      string
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
	s.BaseUri = fmt.Sprintf("http://localhost:%d/system/karas/", kmConfig.System.FrontendPort)
	for _, v := range kmConfig.System.Repositories {
		if len(ctx.StringSlice("repo")) != 0 && !slices.Contains(ctx.StringSlice("repo"), v.Name) {
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
	if len(s.Repositories) == 0 {
		logrus.WithFields(logrus.Fields{
			"any-directories-from": ctx.StringSlice("repo"),
		}).Error("No repository found with the given names")
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

func (s *GitSetup) Run(ctx context.Context) error {
	var pr printer.Printer
	if s.OutputJson {
		pr = printer.NewJsonPrinter()
	} else {
		pr = printer.NewTxtPrinter(s.Hyperlink, s.Color, s.BaseUri)
	}

	// "modified" flag: it is set to true by any goroutine if a file is found
	// this is safe for concurrent use because when the value is updated, it is always to `true`
	modified := false

	nbNotGit := atomic.Uint32{}
	nbHasErr := atomic.Uint32{}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgRepos := sync.WaitGroup{}
	defer func() {
		wgRepos.Wait()
		if nbNotGit.Load() == uint32(len(s.Repositories)) {
			logrus.WithFields(logrus.Fields{
				"num-repositories-not-git": len(s.Repositories),
			}).WithError(ErrNotAGitRepo).Error("None of your repositories are git repositories.")
			return
		}
		if !modified {
			if nbHasErr.Load() == nbNotGit.Load() {
				if uint32(len(s.Repositories))-nbNotGit.Load() == 1 {
					logrus.WithFields(logrus.Fields{
						"repository": s.Repositories[0].Name,
					}).Info("Your git repository is clean. No karaoke to probe.")
				} else {
					logrus.WithFields(logrus.Fields{
						"num-total-git-repositories": uint32(len(s.Repositories)) - nbNotGit.Load(),
					}).Info("All git repositories are clean. No karaoke to probe.")
				}
			} else if nbHasErr.Load() < uint32(len(s.Repositories)) { // at least one repository is clean
				logrus.WithFields(logrus.Fields{
					"num-total-git-repositories":  uint32(len(s.Repositories)) - nbNotGit.Load(),
					"num-failed-git-repositories": nbHasErr.Load() - nbNotGit.Load(),
				}).Info("All other git repositories are clean. No karaoke to probe.")
			}
		}
	}()
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
					nbHasErr.Add(1)
					if errors.Is(err, ErrNotAGitRepo) {
						nbNotGit.Add(1)
					} else if errors.Is(err, ErrUnresolvedMergeConflict) {
						logrus.WithFields(logrus.Fields{
							"repository": repo.Name,
						}).WithError(err).Error("Merge conflict need to be resolved.")
					} else if errors.Is(err, ErrParseError) {
						logrus.WithFields(logrus.Fields{
							"repository": repo.Name,
						}).WithError(err).Error("Unexpected tokens in the output of porcelain `git status` command.")
					} else {
						logrus.WithFields(logrus.Fields{
							"repository": repo.Name,
						}).WithError(err).Error("Failed to read output of porcelain `git status` command.")
					}
					return
				}
				if len(kara) > 0 {
					modified = true
				}
				for _, u := range kara {
					p := filepath.Join(repo.BaseDir, "karaokes", u.String()+".kara.json")
					s.StartWork()
					wg.Go(func() {
						defer s.StopWork()
						app.RunOnFile(ctx, &repo, p, pr)
					})
				}
			})
		}
	}
	return nil
}
