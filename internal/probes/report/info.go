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

// When the test is only to display some infos
func Info(v bool) *report {
	r := reportPool.Get().(*report)
	r.status = status.Info
	r.severity = severity.Info
	if v {
		r.result = result.Passed
	} else {
		r.result = result.Failed
	}
	return r
}
