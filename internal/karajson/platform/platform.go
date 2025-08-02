// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package platform

import (
	"github.com/gofrs/uuid"
)

var (
	// platforms
	Arcade              = uuid.Must(uuid.FromString("bb16f9bd-e5f1-49db-a4ce-dc12cc0fe249"))
	Dreamcast           = uuid.Must(uuid.FromString("cf4a3dd9-affa-49b4-ad4c-67cbadf09c8d"))
	GameCube            = uuid.Must(uuid.FromString("9c2e70a6-ac40-4701-9f3e-268cc52500c8"))
	Mobage              = uuid.Must(uuid.FromString("213a1a15-e73a-4a08-b3f6-b6e2745a82e5"))
	Nintendo3DS         = uuid.Must(uuid.FromString("d7393757-7de7-4247-ba59-120b3ceea9bf"))
	Nintendo64          = uuid.Must(uuid.FromString("17ca27a4-c5f2-425f-ad14-2e2319209129"))
	NintendoDS          = uuid.Must(uuid.FromString("fafc313d-8d8a-4445-90b5-aad7928182ab"))
	NintendoSwitch      = uuid.Must(uuid.FromString("0a880bfc-d0e5-448b-98d2-4ba7b3106641"))
	PC                  = uuid.Must(uuid.FromString("3c15f758-e1bc-4266-8bb9-3c14112f504a"))
	PlayStation         = uuid.Must(uuid.FromString("f84b68ac-943a-495a-9c6e-d90ec959b19e"))
	PlayStation2        = uuid.Must(uuid.FromString("2173004e-21da-45c9-9dfb-26014423c566"))
	PlayStation3        = uuid.Must(uuid.FromString("7f63ff91-7961-4488-888d-4bf46eea3a2a"))
	PlayStation4        = uuid.Must(uuid.FromString("248d92ee-8c02-43bc-86e2-cc6029883564"))
	PlayStation5        = uuid.Must(uuid.FromString("d22db878-440d-46c1-a749-ae575c78983d"))
	PlayStationPortable = uuid.Must(uuid.FromString("6b101aff-4b7b-4b0d-b35e-0f1b67e38c9b"))
	PlayStationVita     = uuid.Must(uuid.FromString("e791f168-068f-4c05-aa36-e8d63726aa06"))
	Saturn              = uuid.Must(uuid.FromString("e6778c89-1d6b-455e-bd53-37a3b12b7140"))
	SegaCD              = uuid.Must(uuid.FromString("b8a3778d-0864-4135-9053-1def02d3262e"))
	SuperNintendo       = uuid.Must(uuid.FromString("b48e6d3d-232e-47e0-be87-e54a1051a764"))
	Wii                 = uuid.Must(uuid.FromString("78d46ef5-c98a-44dd-ad00-5d8c44b6271b"))
	WiiU                = uuid.Must(uuid.FromString("293e4786-9ac8-4c83-84f2-2e6cd14b5e8d"))
	XBOX360             = uuid.Must(uuid.FromString("d56f8d9d-9ccd-4ac8-a164-da0e62d1c060"))
	XboxOne             = uuid.Must(uuid.FromString("22840d76-7030-4c3f-ab35-4a863adce093"))
	XboxSeries          = uuid.Must(uuid.FromString("569fdd6b-1a2c-424e-aab8-93c5891db245"))

	// platforms alias
	DC   = Dreamcast
	GC   = GameCube
	N64  = Nintendo64
	NDS  = NintendoDS
	N3DS = Nintendo3DS
	MOB  = Mobage
	PS2  = PlayStation2
	PS3  = PlayStation3
	PS4  = PlayStation4
	PS5  = PlayStation5
	PSP  = PlayStationPortable
	PSV  = PlayStationVita
	PSX  = PlayStation
	SEGA = SegaCD
)
