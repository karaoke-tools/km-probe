// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/ass"
)

type Probe struct {
	LyricsFile *ass.Ass
	Report     *Report
}

func NewProbe(kara_displayname string, lyricsFile *ass.Ass) *Probe {
	return &Probe{
		LyricsFile: lyricsFile,
		Report:     NewReport(kara_displayname),
	}
}

func (p *Probe) Run(ctx context.Context) error {
	if err := p.CheckAutomation(ctx); err != nil {
		return err
	}
	if err := p.CheckResolution(ctx); err != nil {
		return err
	}
	if err := p.CheckStyle(ctx); err != nil {
		return err
	}
	return nil
}
