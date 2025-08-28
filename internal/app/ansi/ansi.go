// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ansi

const (
	// colors
	Black  = "\033[0;30m"
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[0;33m"
	Blue   = "\033[0;34m"
	Purple = "\033[0;35m"
	Cyan   = "\033[0;36m"
	White  = "\033[0;37m"

	// modifiers
	Bold      = "\033[1m"
	Underline = "\033[4m"

	// reset
	Reset = "\033[0m"
)

func Link(anchor string, display string) string {
	return "\x1B]8;;" + anchor + "\x1B\x5C" + display + "\x1B]8;;\x1B\x5C"
}
