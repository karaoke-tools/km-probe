// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karaokes

import (
	"context"

	// TODO: use go 1.25, <https://github.com/louisroyer/km-probe/issues/16>
	"github.com/louisroyer/km-probe/internal/backport/sync"

	"github.com/louisroyer/km-probe/internal/app"
	"github.com/louisroyer/km-probe/internal/app/printer"
)

// Maximum number of karaokes processed simultaneously.
// It is not useful to increase this number
// because we are bound by the speed of the json encoder.
// Increasing the number of worker will consume more memory
// because we cannot recycle structures when they are still used;
// making the work of the garbage collector more difficult,
// which will slow down everything, and may make the interface irresponsible
// for enough time to be noticable (~1s).
const MAX_WORKERS = 0xFF

// Run on all karaokes of all repositories
func (s *KaraokeSetup) RunAll(ctx context.Context) error {
	var pr printer.Printer
	if s.OutputJson {
		pr = printer.NewJsonPrinter()
	} else {
		pr = printer.NewTxtPrinter(s.Hyperlink, s.Color, s.BaseUri)
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wgRepos := sync.WaitGroup{}
	defer wgRepos.Wait()
	workers := make(chan struct{}, MAX_WORKERS) // maximum number of simultaneous workers
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wgRepos.Go(func() {
				repo.WalkKaraokes(ctx,
					func(ctx context.Context, r *app.Repository, p string) error {
						workers <- struct{}{}
						wg.Go(func() {
							app.RunOnFile(ctx, r, p, pr)
							<-workers
						})
						return nil
					},
				)
			})
		}
	}
	return nil
}
