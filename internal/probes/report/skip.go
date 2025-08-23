// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/report/status"
)

// When the test is not relevant
func Skip(message string) *report {
	r := reportPool.Get().(*report)
	r.status = status.Skipped
	r.severity = severity.Info
	r.message = message // indicate why the test has been skipped
	return r
}
