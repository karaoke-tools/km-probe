// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"os"

	"github.com/louisroyer/km-probe/internal/ass"
)

type SetupParseAss struct {
	File string
}

func NewSetupParseAss(file string) *SetupParseAss {
	return &SetupParseAss{
		File: file,
	}
}

func (s *SetupParseAss) Run(ctx context.Context) error {
	f, err := os.OpenFile(s.File, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = ass.Parse(ctx, f)
	if err != nil {
		return err
	}
	return nil
}
