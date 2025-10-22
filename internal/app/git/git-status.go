// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package git

import (
	"bufio"
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrUnresolvedMergeConflict = errors.New("Unresolved merge conflict")
	ErrNotAGitRepo             = errors.New("Not a git repository")
	ErrParseError              = errors.New("Error while parsing git status output")
)

func GitModifiedKaras(ctx context.Context, path string) ([]uuid.UUID, error) {
	kara := []uuid.UUID{}
	// verify this is a git repository, to avoid doing all the below for nothing
	// because we can only check the result of the command after parsing stdout
	if _, err := os.Stat(filepath.Join(path, ".git")); errors.Is(err, fs.ErrNotExist) {
		return kara, ErrNotAGitRepo
	}
	cmd := exec.CommandContext(ctx, "git", "-C", filepath.Clean(path), "status", "--porcelain=v2")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return kara, err
	}
	if err := cmd.Start(); err != nil {
		return kara, err
	}
	scanner := bufio.NewScanner(stdout)
	var errScan error = nil
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return []uuid.UUID{}, ctx.Err()
		default:
			if errScan != nil {
				continue
			}
			line := string(scanner.Bytes())
			file := ""
			if strings.HasPrefix(line, "#") {
				// header
				continue
			}
			c, after, ok := strings.Cut(line, " ")
			if !ok {
				errScan = ErrParseError
				continue
			}
			switch c {
			case "?":
				// untracked item
				file = after
			case "1":
				// ordinary changed entry
				xy, _, ok := strings.Cut(after, " ")
				if !ok {
					errScan = ErrParseError
					continue
				}
				if len(xy) != 2 {
					errScan = ErrParseError
					continue
				}
				if string(xy[1]) == "D" {
					// deleted file
					continue
				}
				sp := strings.SplitN(after, " ", 8)
				if len(sp) != 8 {
					errScan = ErrParseError
					continue
				}
				file = sp[7]
			case "2":
				// renamed or copied entry
				sp := strings.SplitN(after, " ", 9)
				if len(sp) != 9 {
					errScan = ErrParseError
					continue
				}
				file, _, ok = strings.Cut(sp[8], "\t")
				if !ok {
					errScan = ErrParseError
					continue
				}
			case "u":
				// unmerged entry
				// -> there is a merge conflict
				errScan = ErrUnresolvedMergeConflict
				continue
			case "!":
				// ignored item
				// -> we ignore it as well
				continue
			default:
				// unexpected symbol
				errScan = ErrParseError
				continue
			}
			file, ok = strings.CutPrefix(file, "karaokes/")
			if ok {
				file, ok = strings.CutSuffix(file, ".kara.json")
			} else {
				file, ok = strings.CutPrefix(file, "lyrics/")
				if ok {
					file, ok = strings.CutSuffix(file, ".ass")
				}
			}
			if !ok {
				// not a karaoke, may be another file
				continue
			}
			k, err := uuid.FromString(file)
			if err != nil {
				// repo not using UUIDs ?
				continue
			}
			if !slices.Contains(kara, k) {
				kara = append(kara, k)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return kara, err
	}
	return kara, errScan
}
