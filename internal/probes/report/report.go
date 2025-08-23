// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"encoding/json"
	"sync"

	"github.com/louisroyer/km-probe/internal/probes/report/result"
	"github.com/louisroyer/km-probe/internal/probes/report/severity"
	"github.com/louisroyer/km-probe/internal/probes/report/status"
)

type Report interface {
	json.Marshaler

	Result() result.Result // true: passed, false: failed
	Status() status.Status // completed, aborted, skipped, etc.
	Severity() severity.Severity
	Message() string

	// Delete should be used when the Report is no longer useful.
	// This allows to recycle the memory for future usage.
	Delete()
}

func (r *report) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status   string `json:"status"`
		Result   string `json:"result"`
		Message  string `json:"message,omitempty"`
		Severity string `json:"severity"`
	}{
		Status:   r.status.String(),
		Result:   r.result.String(),
		Message:  r.message,
		Severity: r.severity.String(),
	})
}

type report struct {
	status   status.Status
	result   result.Result
	message  string
	severity severity.Severity

	// when not set, the pool is not used
	recycleAfterUse bool
}

// pool to recycle memory
var reportPool = sync.Pool{
	New: func() any {
		return &report{
			recycleAfterUse: true,
		}
	},
}

// Delete should be used when the report is no longer useful.
// This allows to recycle the memory for future usage.
func (r *report) Delete() {
	if !r.recycleAfterUse {
		return
	}
	r.severity = severity.Unknown
	r.status = status.Unknown
	r.result = result.Unknown
	r.message = ""
	reportPool.Put(r)
}

func (r *report) Status() status.Status {
	return r.status
}

func (r *report) Result() result.Result {
	return r.result
}

func (r *report) Severity() severity.Severity {
	return r.severity
}

func (r *report) Message() string {
	return r.message
}
