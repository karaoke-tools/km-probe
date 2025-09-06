// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karadata

import (
	"context"
	"os"
	"path"

	"github.com/karaoke-tools/km-probe/internal/ass"
	"github.com/karaoke-tools/km-probe/internal/karajson"
)

// Karaoke information
type KaraData struct {
	KaraJson *karajson.KaraJson // metadata of the karaoke
	Lyrics   []*ass.Ass         // lyrics of the karaoke
}

// Create a new `KaraData` from a `KaraJson`
func FromKaraJson(ctx context.Context, basedir string, karaJson *karajson.KaraJson) (*KaraData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if len(karaJson.Medias) == 0 {
			return nil, ErrNoMedias
		}
		data := KaraData{
			KaraJson: karaJson,
			Lyrics:   make([]*ass.Ass, 0, len(karaJson.Medias[0].Lyrics)),
		}
		// TODO: update this when multi-track drifting is released
		for _, l := range data.KaraJson.Medias[0].Lyrics {
			lyricsPath := path.Join(basedir, "lyrics", l.Filename)
			f, err := os.OpenFile(lyricsPath, os.O_RDONLY, 0)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			lyrics, err := ass.Parse(ctx, f)
			if err != nil {
				return nil, err
			}
			data.Lyrics = append(data.Lyrics, lyrics)
		}
		return &data, nil
	}
}
