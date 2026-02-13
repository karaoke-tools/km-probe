// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package genre

import (
	"github.com/gofrs/uuid/v5"
)

var (
	// genres
	BoysLove         = uuid.Must(uuid.FromString("724d2d80-a531-4579-bdaf-4a0db9a313d9"))
	GirlsLove        = uuid.Must(uuid.FromString("76378421-153e-410d-a62c-2db43d5ed5fd"))
	Idol             = uuid.Must(uuid.FromString("caa86df2-0d59-474b-885c-f240a9e891b0"))
	Isekai           = uuid.Must(uuid.FromString("4b80eca6-f44f-4262-a8d0-b6879d6b70cb"))
	MagicalGirl      = uuid.Must(uuid.FromString("b84e28f9-db1e-447f-b339-7954c3592523"))
	Mecha            = uuid.Must(uuid.FromString("f60ba57d-5ef4-49e0-b93b-d0dcbcff6592"))
	Otome            = uuid.Must(uuid.FromString("7025ae3c-7bd6-4787-8912-2c15063343bb"))
	Shoujo           = uuid.Must(uuid.FromString("8b6ace6f-a59e-4740-b3ff-e1618720383c"))
	Shounen          = uuid.Must(uuid.FromString("a0aeef4a-6428-45ff-a6e1-468b595930c2"))
	Tokusatsu        = uuid.Must(uuid.FromString("d3fac9ab-630d-402d-a392-3d2450c3e62e"))
	VocalSynthesizer = uuid.Must(uuid.FromString("45623cae-6d68-4304-a49b-896d1d6f4580"))

	// genres alias
	BL       = BoysLove
	GL       = GirlsLove
	MAGIC    = MagicalGirl
	OTO      = Otome
	SYNTH    = VocalSynthesizer
	TKU      = Tokusatsu
	Vocaloid = VocalSynthesizer
)
