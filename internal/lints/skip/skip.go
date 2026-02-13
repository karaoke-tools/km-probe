// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package skip

import (
	"context"

	"github.com/karaoke-tools/km-probe/internal/karadata"
)

type Condition interface {
	Result(context.Context, *karadata.KaraData) (bool, string, error)
}
