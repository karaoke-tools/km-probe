// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karaokes

import (
	"context"
	"slices"

	"github.com/louisroyer/km-probe/internal/app"
	"github.com/louisroyer/km-probe/internal/app/setup"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type KaraokeSetup struct {
	*setup.Setup
	Repositories []app.Repository
	Uuids        []uuid.UUID
}

func FromCli(ctx *cli.Context) (*KaraokeSetup, error) {
	s := &KaraokeSetup{
		Setup:        setup.FromCli(ctx),
		Repositories: make([]app.Repository, 0),
		Uuids:        make([]uuid.UUID, 0),
	}

	// parse uuid
	enabledUuids := ctx.StringSlice("uuid")
	for _, enabledUuid := range enabledUuids {
		if u, err := uuid.FromString(enabledUuid); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"uuid-argument": enabledUuid,
			}).Error("Could not parse uuid")
			return nil, err
		} else {
			s.Uuids = append(s.Uuids, u)
		}
	}

	kmConfig, err := app.LoadConf()
	if err != nil {
		return nil, err
	}
	for _, v := range kmConfig.System.Repositories {
		if len(ctx.StringSlice("repository")) != 0 && !slices.Contains(ctx.StringSlice("repository"), v.Name) {
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
			"any-directories-from": ctx.StringSlice("repository"),
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

func (s *KaraokeSetup) Run(ctx context.Context) error {
	if len(s.Uuids) == 0 {
		return s.RunAll(ctx)
	}
	return s.RunByUuid(ctx)
}
