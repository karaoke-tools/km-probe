// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package language

import (
	"github.com/gofrs/uuid"
)

var (
	// langs (only a subset of all languages)
	JPN = uuid.Must(uuid.FromString("4dcf9614-7914-42aa-99f4-dbce2e059133")) // japanese
	ENG = uuid.Must(uuid.FromString("de5eda1c-5fb3-46a6-9606-d4554fc5a1d6")) // english
	MUL = uuid.Must(uuid.FromString("1a23082a-620c-4bd8-9930-64d1798515c1")) // multiple languages
	ZXX = uuid.Must(uuid.FromString("51e68bf8-01ab-4210-8e1f-26cbf2bd68be")) // no linguistic content
)
