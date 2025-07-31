// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package karajson

import (
	"github.com/gofrs/uuid"
)

var (
	// songtype
	TypeAudioOnly = uuid.Must(uuid.FromString("42a262ae-acba-4ab5-a446-c5789c96c821"))

	// misc
	GroupSinging = uuid.Must(uuid.FromString("1d6ceadf-a885-47e5-b23b-c07bf668a424"))
	Unavailable  = uuid.Must(uuid.FromString("7a1ad419-4001-484a-bd9c-d3e04bf54529"))

	// collection
	CollectionAsia      = uuid.Must(uuid.FromString("dbcf2c22-524d-4708-99bb-601703633927"))
	CollectionGeekOtaku = uuid.Must(uuid.FromString("c7db86a0-ff64-4044-9be4-66dd1ef1d1c1"))
	CollectionShitpost  = uuid.Must(uuid.FromString("2fa2fe3f-bb56-45ee-aa38-eae60e76f224"))
	CollectionKana      = uuid.Must(uuid.FromString("f2462778-f986-4844-a4b8-e1d3ccdb861b"))
	CollectionWest      = uuid.Must(uuid.FromString("efe171c0-e8a1-4d03-98c0-60ecf741ad52"))

	// langs
	LangJPN = uuid.Must(uuid.FromString("4dcf9614-7914-42aa-99f4-dbce2e059133"))
	LangENG = uuid.Must(uuid.FromString("de5eda1c-5fb3-46a6-9606-d4554fc5a1d6"))

	// versions
	VersionOffVocal = uuid.Must(uuid.FromString("c0cc87b9-55b9-40f0-878a-fbb9e34c151e"))

	// warnings
	WarningR18Lyrics = uuid.Must(uuid.FromString("e2b8419f-1d5a-44ad-a62c-d7765493190d"))
	WarningR18Media  = uuid.Must(uuid.FromString("e82ce681-6d7b-4fb6-abe4-daa8aaa9bbf9"))
	WarningSpoiler   = uuid.Must(uuid.FromString("24371984-5e4c-4485-a937-fb0c480ca23b"))
	WarningEpilepsy  = uuid.Must(uuid.FromString("51288600-29e0-4e41-a42b-77f0498e5691"))

	// years
	Years1950 = uuid.Must(uuid.FromString("de0c8e96-9934-4251-9142-b7a5ee70dd6f"))
	Years1960 = uuid.Must(uuid.FromString("b4240f1f-200d-4d10-b5eb-d3a740eed323"))
	Years1970 = uuid.Must(uuid.FromString("57e5e0c0-e02e-4871-a6cb-d48ea604ebab"))
	Years1980 = uuid.Must(uuid.FromString("7792cc27-9ff9-4d47-9287-2338b7db1575"))
	Years1990 = uuid.Must(uuid.FromString("24db50aa-a126-4d34-a48a-e45db79c4245"))
	Years2000 = uuid.Must(uuid.FromString("3c14c864-533c-49e6-a3cc-6969f1233775"))
	Years2010 = uuid.Must(uuid.FromString("ae613721-1fbe-480d-ba6b-d6d0702b184d"))
	Years2020 = uuid.Must(uuid.FromString("8e209ae4-4381-438f-89dc-1b21752027d2"))
)
