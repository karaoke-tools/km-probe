// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karajson

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gofrs/uuid"
)

func FromFile(path string) (*KaraJson, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	kara := new(KaraJson)
	if err := json.Unmarshal(content, kara); err != nil {
		return nil, err
	}
	return kara, nil
}

type KaraJson struct {
	Header Header  `json:"header"`
	Medias []Media `json:"medias"`
	Data   Data    `json:"data"`
}

type Header struct {
	Version     uint32 `json:"version"`
	Description string `json:"description"`
}

type Media struct {
	Default  bool    `json:"default"`
	Duration uint32  `json:"duration"`
	Filename string  `json:"filename"`
	Filesize uint64  `json:"filesize"`
	Loudnorm string  `json:"loudnorm"`
	Lyrics   []Lyric `json:"lyrics"`
	Version  string  `json:"version"`
}

type Lyric struct {
	Default  bool   `json:"default"`
	Filename string `json:"filename"`
	Version  string `json:"version"`
}

type Data struct {
	CreatedAt             time.Time         `json:"created_at"`
	IgnoreHooks           bool              `json:"ignore-hooks"`
	Kid                   uuid.UUID         `json:"kid"`
	ModifiedAt            time.Time         `json:"modified_at"`
	Parents               []uuid.UUID       `json:"parents"`
	Repository            string            `json:"repository"`
	Songname              string            `json:"songname"`
	Tags                  Tags              `json:"tags"`
	Titles                map[string]string `json:"titles"`
	TitlesAliases         []string          `json:"titles_aliases"`
	TitlesDefaultLanguage string            `json:"titles_default_language"`
	Year                  json.Number       `json:"year"`
}

type Tags struct {
	Authors      []uuid.UUID `json:"authors"`
	Collections  []uuid.UUID `json:"collections"`
	Creators     []uuid.UUID `json:"creators"`
	Families     []uuid.UUID `json:"families"`
	Groups       []uuid.UUID `json:"groups"`
	Langs        []uuid.UUID `json:"langs"`
	Misc         []uuid.UUID `json:"misc"`
	Origins      []uuid.UUID `json:"origins"`
	Platforms    []uuid.UUID `json:"platforms"`
	Series       []uuid.UUID `json:"series"`
	Singers      []uuid.UUID `json:"singers"`
	Singergroups []uuid.UUID `json:"singergroups"`
	Songtypes    []uuid.UUID `json:"songtypes"`
	Songwriters  []uuid.UUID `json:"songwriters"`
	Versions     []uuid.UUID `json:"versions"`
	Warnings     []uuid.UUID `json:"warnings"`
}
