// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"errors"
	"io/fs"

	"github.com/karaoke-tools/km-probe/internal/app/printer"
	"github.com/karaoke-tools/km-probe/internal/karadata"
	"github.com/karaoke-tools/km-probe/internal/karajson"

	"github.com/sirupsen/logrus"
)

// Parse the karaoke, run probes, and display result
// p is the filepath to the .kara.json file
func RunOnFile(ctx context.Context, repo *Repository, p string, pr printer.Printer) error {
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
	}
	aggregator := pr.Aggregator()
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
	if err := pr.Encode(ctx, aggregator); err != nil {
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
