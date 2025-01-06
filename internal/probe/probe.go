// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"os"
	"path"

	"github.com/louisroyer/km-probe/internal/ass"
	"github.com/louisroyer/km-probe/internal/karajson"
)

type Probe struct {
	KaraJson *karajson.KaraJson
	Lyrics   *ass.Ass
	Report   *Report
}

func FromKaraJson(ctx context.Context, basedir string, karaJson *karajson.KaraJson) (*Probe, error) {
	probe := Probe{
		KaraJson: karaJson,
		Report:   NewReport(karaJson.Data.Songname),
	}
	if len(karaJson.Medias[0].Lyrics) == 0 {
		return nil, ErrNoLyrics
	}
	lyricsPath := path.Join(basedir, "lyrics", karaJson.Medias[0].Lyrics[0].Filename)
	f, err := os.OpenFile(lyricsPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lyrics, err := ass.Parse(ctx, f)
	if err != nil {
		return nil, err
	}
	probe.Lyrics = lyrics
	return &probe, nil
}

func (p *Probe) Run(ctx context.Context) error {
	if err := p.checkLiveDownloadProbablyAllowed(ctx); err != nil {
		return err
	}
	if err := p.checkAutomation(ctx); err != nil {
		return err
	}
	if err := p.checkResolution(ctx); err != nil {
		return err
	}
	if err := p.checkStyle(ctx); err != nil {
		return err
	}
	return nil
}
