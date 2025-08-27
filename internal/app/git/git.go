// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"errors"

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

func (s *GitSetup) Run(ctx context.Context) error {
	return errors.New("Not implemented")
}
