// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karajson

import (
	"github.com/gofrs/uuid"
)

var (
	// misc
	GroupSinging = uuid.Must(uuid.FromString("1d6ceadf-a885-47e5-b23b-c07bf668a424"))
	Unavailable  = uuid.Must(uuid.FromString("7a1ad419-4001-484a-bd9c-d3e04bf54529"))

	// collection
	CollectionAsia      = uuid.Must(uuid.FromString("dbcf2c22-524d-4708-99bb-601703633927"))
	CollectionGeekOtaku = uuid.Must(uuid.FromString("c7db86a0-ff64-4044-9be4-66dd1ef1d1c1"))
	CollectionShitpost  = uuid.Must(uuid.FromString("2fa2fe3f-bb56-45ee-aa38-eae60e76f224"))
	CollectionKana      = uuid.Must(uuid.FromString("f2462778-f986-4844-a4b8-e1d3ccdb861b"))
	CollectionWest      = uuid.Must(uuid.FromString("efe171c0-e8a1-4d03-98c0-60ecf741ad52"))
)
