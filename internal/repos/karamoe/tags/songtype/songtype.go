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
	AnimeAmateurMusicVideo = uuid.Must(uuid.FromString("a6c79ce5-89ee-4d50-afe8-3abd7317f6c2")) // AMV
	AudioOnly              = uuid.Must(uuid.FromString("42a262ae-acba-4ab5-a446-c5789c96c821")) // AUDIO
	CharacterSong          = uuid.Must(uuid.FromString("a8a2616c-b5a1-4df3-ad6f-89291bc67dee")) // CS
	Commercial             = uuid.Must(uuid.FromString("2ddb5358-e674-46fa-a6e1-7f5c5d56f8fa")) // CM
	Concert                = uuid.Must(uuid.FromString("a0167949-580c-4de3-bf13-497e462e02f3")) // LIVE
	Ending                 = uuid.Must(uuid.FromString("38c77c56-2b95-4040-b676-0994a8cb0597")) // ED
	InsertSong             = uuid.Must(uuid.FromString("5e5250d9-351a-4a82-98eb-55db50ad8962")) // IN
	ImageSong              = uuid.Must(uuid.FromString("10a1ad3e-a05c-4f5c-84b6-f491e3e3a92e")) // IS
	MusicVideo             = uuid.Must(uuid.FromString("7be1b15c-cff8-4b37-a649-5c90f3d569a9")) // MV
	Opening                = uuid.Must(uuid.FromString("f02ad9b3-0bd9-4aad-85b3-9976739ba0e4")) // OP
	OtherSong              = uuid.Must(uuid.FromString("97769615-a2e5-4f36-8c23-b2ce2ce3c460")) // OT
	PromotionalVideo       = uuid.Must(uuid.FromString("66e195b9-64ea-432d-a917-2e2c88d26052")) // PV
	Streaming              = uuid.Must(uuid.FromString("55ce3d79-dcc2-453c-b00a-60ce0c1eba1c")) // STREAM

	// songtype alias
	AMV    = AnimeAmateurMusicVideo
	AUDIO  = AudioOnly
	CM     = Commercial
	CS     = CharacterSong
	ED     = Ending
	IN     = InsertSong
	IS     = ImageSong
	LIVE   = Concert
	MV     = MusicVideo
	OP     = Opening
	OT     = OtherSong
	PV     = PromotionalVideo
	STREAM = Streaming
)
