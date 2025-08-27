// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karaokes

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/app"

	"github.com/sirupsen/logrus"
)

func (s *KaraokeSetup) RunSingle(ctx context.Context) (err error) {
	printer := app.NewPrinter()

	// found flag: it is set to true by any goroutine if a file is found
	// this is safe for concurrent use because when the value is updated, it is always to `true`
	found := false

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		if !found {
			// `err` is a named return value: this allow us to modify it inside the defer
			err = app.ErrKaraokeNotFound
			logrus.WithFields(logrus.Fields{
				"uuid": s.Uuid,
			}).WithError(err).Error("Karaoke not found")
		}
	}()

	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wg.Go(func() {
				fp := filepath.Join(repo.BaseDir, "karaokes", s.Uuid.String()+".kara.json")
				err := app.RunOnFile(ctx, &repo, fp, printer)
				if err == nil || !errors.Is(err, fs.ErrNotExist) {
					found = true
				}
			})
		}
	}
	return nil
}
