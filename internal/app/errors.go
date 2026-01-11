// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"errors"
)

var (
	ErrConfigNotFound         = errors.New("KM configuration file not found")
	ErrDataRepositoryNotFound = errors.New("KM data repository not found")
	ErrKaraokeNotFound        = errors.New("karaoke not found")
	ErrDuplicateKaraoke       = errors.New("karaoke found in multiple repositories")
)
