// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package severity

type Severity int

const (
	Unknown  Severity = iota // no info
	Critical                 // there is something to be corrected
	Warning                  // maintainer should have a look at it, but may decide to ignore it
	Info                     // maintainer can safely ignore this info
)

func (s Severity) String() string {
	switch s {
	case Critical:
		return "critical"
	case Warning:
		return "warning"
	case Info:
		return "info"
	case Unknown:
		fallthrough
	default:
		return "unknown"
	}
}
