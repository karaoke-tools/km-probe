// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package result

type Result int

const (
	Unknown Result = iota
	Passed
	Failed
)

func (r Result) String() string {
	switch r {
	case Passed:
		return "passed"
	case Failed:
		return "failed"
	case Unknown:
		fallthrough
	default:
		return "unknown status"
	}
}
