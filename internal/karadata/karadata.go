// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karadata

import (
	"context"
	"os"
	"path"

	"github.com/louisroyer/km-probe/internal/ass"
	"github.com/louisroyer/km-probe/internal/karajson"
)

type KaraData struct {
	KaraJson *karajson.KaraJson
	Lyrics   *ass.Ass
}

func FromKaraJson(ctx context.Context, basedir string, karaJson *karajson.KaraJson) (*KaraData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		data := KaraData{
			KaraJson: karaJson,
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
		data.Lyrics = lyrics
		return &data, nil
	}
}
