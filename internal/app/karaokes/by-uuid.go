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
	"github.com/louisroyer/km-probe/internal/app/printer"

	"github.com/sirupsen/logrus"
)

func (s *KaraokeSetup) RunByUuid(ctx context.Context) (err error) {
	var pr printer.Printer
	if s.OutputJson {
		pr = printer.NewJsonPrinter()
	} else {
		pr = printer.NewTxtPrinter(s.Hyperlink, s.Color, s.BaseUri)
	}

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
				"any-uuids-from": s.Uuids,
			}).WithError(err).Error("No karaoke found matching given criterias")
		}
	}()

	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for _, u := range s.Uuids {
				wg.Go(func() {
					fp := filepath.Join(repo.BaseDir, "karaokes", u.String()+".kara.json")
					err := app.RunOnFile(ctx, &repo, fp, pr)
					if err == nil || !errors.Is(err, fs.ErrNotExist) {
						found = true
					}
				})
			}
		}
	}
	return nil
}
