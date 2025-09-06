// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"github.com/karaoke-tools/km-probe/internal/probes/report/status"
)

// `abort` has no custom information, and can be reused
// without re-allocation or cleaning
var abort = report{
	status: status.Aborted,
}

// When test has been canceled
func Abort() *report {
	return &abort
}
