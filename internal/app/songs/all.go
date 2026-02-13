// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package songs

import (
	"context"
	"sync"

	"github.com/karaoke-tools/km-probe/internal/app"
	"github.com/karaoke-tools/km-probe/internal/app/printer"
)

// Run on all songs of all repositories
func (s *SongsSetup) RunAll(ctx context.Context) error {
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
	for _, repo := range s.Repositories {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wgRepos.Go(func() {
				repo.WalkSongs(ctx,
					func(ctx context.Context, r *app.Repository, p string) error {

						s.StartWork()
						wg.Go(func() {
							defer s.StopWork()
							app.RunOnFile(ctx, r, p, pr)
						})
						return nil
					},
				)
			})
		}
	}
	return nil
}
