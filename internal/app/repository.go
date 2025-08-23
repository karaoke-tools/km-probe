// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
)

type Repository struct {
	Name      string
	BaseDir   string
	MediaPath string
}

// Run the given f function on all karaokes of a repository.
func (repo *Repository) WalkKaraokes(ctx context.Context, f func(ctx context.Context, repo *Repository, p string) error) error {
	return filepath.WalkDir(filepath.Join(repo.BaseDir, "karaokes"), func(p string, d fs.DirEntry, err error) error {
		// check file metadata
		select {
		case <-ctx.Done():
			return filepath.SkipDir
		default:
			if d.IsDir() {
				if p == filepath.Join(repo.BaseDir, "karaokes") {
					return nil
				}
				return filepath.SkipDir
			}
			if !strings.HasSuffix(d.Name(), ".kara.json") {
				return nil
			}
		}
		return f(ctx, repo, p)
	})
}
