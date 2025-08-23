// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package colour

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidFormat = errors.New("Invalid colour format")
)

type ColourComponent uint8

type Colour struct {
	Alpha ColourComponent
	Blue  ColourComponent
	Green ColourComponent
	Red   ColourComponent
}

func (c Colour) String() string {
	return fmt.Sprintf("&H%02X%02X%02X%02X", c.Alpha, c.Blue, c.Green, c.Red)
}

func FromString(s string) (Colour, error) {
	if len(s) != 10 {
		return Colour{}, ErrInvalidFormat
	}
	if !strings.HasPrefix(s, "&H") {
		return Colour{}, ErrInvalidFormat
	}
	if decoded, err := hex.DecodeString(s[2:]); err == nil && len(decoded) == 4 {
		return Colour{
			Alpha: ColourComponent(decoded[0]),
			Blue:  ColourComponent(decoded[1]),
			Green: ColourComponent(decoded[2]),
			Red:   ColourComponent(decoded[3]),
		}, nil
	}
	return Colour{}, ErrInvalidFormat
}

var (
	Black = Colour{}
	White = Colour{Blue: 0xFF, Green: 0xFF, Red: 0xFF}
)
