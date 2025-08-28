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
	"sync/atomic"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/app"
	"github.com/louisroyer/km-probe/internal/app/printer"

	"github.com/sirupsen/logrus"
)

func (s *KaraokeSetup) RunByUuid(ctx context.Context) error {
	var pr printer.Printer
	if s.OutputJson {
		pr = printer.NewJsonPrinter()
	} else {
		pr = printer.NewTxtPrinter(s.Hyperlink, s.Color, s.BaseUri)
	}

	nbFound := make([]atomic.Uint32, len(s.Uuids))

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		for i, n := range nbFound {
			if n.Load() == 0 {
				logrus.WithFields(logrus.Fields{
					"uuid": s.Uuids[i],
				}).WithError(app.ErrKaraokeNotFound).Error("No karaoke not found with this UUID.")
			} else if n.Load() > 1 {
				logrus.WithFields(logrus.Fields{
					"uuid":     s.Uuids[i],
					"nb-found": n.Load(),
				}).WithError(app.ErrDuplicateKaraoke).Error("Found multiple karaokes with this UUID (in multiple repositories).")
			}
		}

	}()

	for i, u := range s.Uuids {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for _, repo := range s.Repositories {
				s.StartWork()
				wg.Go(func() {
					defer s.StopWork()
					fp := filepath.Join(repo.BaseDir, "karaokes", u.String()+".kara.json")
					err := app.RunOnFile(ctx, &repo, fp, pr)
					if err == nil || !errors.Is(err, fs.ErrNotExist) {
						nbFound[i].Add(1)
					}
				})
			}
		}
	}
	return nil
}
