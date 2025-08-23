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
	Result() result.Result // true: passed, false: failed
	Status() status.Status // completed, aborted, skipped, etc.
	Severity() severity.Severity
	Message() string
	json.Marshaler
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
	recycleAfterUse bool
	status          status.Status
	result          result.Result
	message         string
	severity        severity.Severity
}

// `pass` has no custom information, and can be reused
// without re-allocation or cleaning
var pass = report{
	status:   status.Completed,
	severity: severity.Info,
	result:   result.Passed,
}

var reportPool = sync.Pool{
	New: func() any {
		return &report{
			recycleAfterUse: true,
		}
	},
}

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

// When the issue is not detected
func Pass() *report {
	return &pass
}

// When the issue is detected
func Fail(severity severity.Severity, message string) *report {
	r := reportPool.Get().(*report)
	r.severity = severity
	r.message = message // indicate what action must be done
	r.status = status.Completed
	r.result = result.Failed
	return r
}

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

// When test has been canceled
func Abort() *report {
	r := reportPool.Get().(*report)
	r.status = status.Aborted
	return r
}

// When the test is not relevant
func Skip(message string) *report {
	r := reportPool.Get().(*report)
	r.status = status.Skipped
	r.severity = severity.Info
	r.message = message // indicate why the test has been skipped
	return r
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
