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
	status string
	result bool
}

func Pass() *report {
	return &report{
		status: "completed",
		result: true,
	}
}

func Fail() *report {
	return &report{
		status: "completed",
		result: false,
	}
}

func Info(v bool) *report {
	return &report{
		status: "info",
		result: v,
	}
}

func Abort() *report {
	return &report{
		status: "aborted",
	}
}

func (r *report) Status() string {
	return r.status
}

func (r *report) Result() bool {
	return r.result
}

func (r *report) String() string {
	if r.status == "info" {
		if r.result {
			return "info: yes"
		} else {
			return "info: no"
		}
	}
	if r.status != "completed" {
		return r.status
	}
	if r.result {
		return "passed"
	} else {
		return "failed"
	}
}
