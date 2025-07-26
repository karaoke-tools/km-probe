// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

type Report interface {
	Result() bool
	Status() string
	String() string
}

type report struct {
	status Status
	result bool
}

// When the issue is not detected
func Pass() *report {
	return &report{
		status: StatusCompleted,
		result: true,
	}
}

// When the issue is detected
func Fail() *report {
	return &report{
		status: StatusCompleted,
		result: false,
	}
}

// When the test is only to display some infos
func Info(v bool) *report {
	return &report{
		status: StatusInfo,
		result: v,
	}
}

// When test has been canceled
func Abort() *report {
	return &report{
		status: StatusAborted,
	}
}

// When the test is not relevant
func Skip() *report {
	return &report{
		result: true,
		status: StatusSkipped,
	}
}

func (r *report) Status() string {
	return r.status.String()
}

func (r *report) Result() bool {
	return r.result
}

func (r *report) String() string {
	if r.status == StatusInfo {
		if r.result {
			return "info: yes"
		} else {
			return "info: no"
		}
	}
	if r.status != StatusCompleted {
		return r.Status()
	}
	if r.result {
		return "passed"
	} else {
		return "failed"
	}
}
