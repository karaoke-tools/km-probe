// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package songtype

import (
	"github.com/gofrs/uuid/v5"
)

var (
	// songtype
	AudioOnly  = uuid.Must(uuid.FromString("42a262ae-acba-4ab5-a446-c5789c96c821")) // AUDIO
	Concert    = uuid.Must(uuid.FromString("a0167949-580c-4de3-bf13-497e462e02f3")) // LIVE
	MusicVideo = uuid.Must(uuid.FromString("7be1b15c-cff8-4b37-a649-5c90f3d569a9")) // MV
	OtherSong  = uuid.Must(uuid.FromString("97769615-a2e5-4f36-8c23-b2ce2ce3c460")) // OT

	// songtype alias
	AUDIO = AudioOnly
	LIVE  = Concert
	MV    = MusicVideo
	OT    = OtherSong
)
