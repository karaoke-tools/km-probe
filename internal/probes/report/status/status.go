// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package status

type Status int

const (
	Completed Status = iota
	Info
	Aborted
	Skipped
)

func (s Status) String() string {
	switch s {
	case Completed:
		return "completed"
	case Info:
		return "info"
	case Aborted:
		return "aborted"
	case Skipped:
		return "skipped"
	default:
		return "unknown status"
	}
}
