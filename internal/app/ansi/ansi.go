// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ansi

const (
	Underline = "\033[4m"
	Reset     = "\033[0m"
)

func Link(anchor string, display string) string {
	return "\x1B]8;;" + anchor + "\x1B\x5C" + display + "\x1B]8;;\x1B\x5C"
}
