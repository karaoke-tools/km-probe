// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

import (
	"encoding/json"

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
}

// When the issue is not detected
func Pass() *report {
	return &report{
		status:   status.Completed,
		severity: severity.Info,
		result:   result.Passed,
	}
}

// When the issue is detected
func Fail(severity severity.Severity, message string) *report {
	return &report{
		severity: severity,
		message:  message, // indicate what action must be done
		status:   status.Completed,
		result:   result.Failed,
	}
}

// When the test is only to display some infos
func Info(v bool) *report {
	var r result.Result
	if v {
		r = result.Passed
	} else {
		r = result.Failed
	}
	return &report{
		status:   status.Info,
		result:   r,
		severity: severity.Info,
	}
}

// When test has been canceled
func Abort() *report {
	return &report{
		status: status.Aborted,
	}
}

// When the test is not relevant
func Skip(message string) *report {
	return &report{
		status:   status.Skipped,
		severity: severity.Info,
		message:  message, // indicate why the test has been skipped
	}
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
