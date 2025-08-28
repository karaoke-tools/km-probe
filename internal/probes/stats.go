// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probes

type Stats struct {

	// status: completed
	Passed         uint32 `json:"passed"`
	FailedCritical uint32 `json:"failed-critical"`
	FailedWarning  uint32 `json:"failed-warning"`
	FailedInfo     uint32 `json:"failed-info"`

	// other status
	Aborted uint32 `json:"aborted"`
	Skipped uint32 `json:"skipped"`
}

func (s *Stats) Reset() {
	s.Passed = 0
	s.FailedCritical = 0
	s.FailedWarning = 0
	s.FailedInfo = 0
	s.Aborted = 0
	s.Skipped = 0
}
