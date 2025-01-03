// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("Probe not implemented")
	ErrMalformedAss   = errors.New("Malformed ASS file")
	ErrMalformedStyle = errors.New("Malformed Style line")
)
