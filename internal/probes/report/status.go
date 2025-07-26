// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package report

type Status int

const (
	StatusCompleted Status = iota
	StatusInfo
	StatusAborted
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusCompleted:
		return "completed"
	case StatusInfo:
		return "info"
	case StatusAborted:
		return "aborted"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown status"
	}
}
