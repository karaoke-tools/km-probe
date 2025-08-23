// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"github.com/louisroyer/km-probe/internal/probes/report/result"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/report/status"
)

// `pass` has no custom information, and can be reused
// without re-allocation or cleaning
var pass = report{
	status:   status.Completed,
	severity: severity.Info,
	result:   result.Passed,
}

// When the issue is not detected
func Pass() *report {
	return &pass
}
