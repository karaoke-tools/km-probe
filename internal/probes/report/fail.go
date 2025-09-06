// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"github.com/karaoke-tools/km-probe/internal/probes/report/result"
	"github.com/karaoke-tools/km-probe/internal/probes/report/severity"
	"github.com/karaoke-tools/km-probe/internal/probes/report/status"
)

// When the issue is detected
func Fail(severity severity.Severity, message string) *report {
	r := reportPool.Get().(*report)
	r.severity = severity
	r.message = message // indicate what action must be done
	r.status = status.Completed
	r.result = result.Failed
	return r
}
