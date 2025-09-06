// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package backports the go 1.25 WaitGroup.Go() feature.
// It can be removed after the toolchain has been updated to go 1.25.
package sync

import (
	"sync"
)

type WaitGroup struct {
	sync.WaitGroup
}
