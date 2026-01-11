// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package colour_test

import (
	"testing"

	"github.com/karaoke-tools/km-probe/internal/ass/style/colour"
)

func TestColourString(t *testing.T) {
	if got := colour.White.String(); got != "&H00FFFFFF" {
		t.Errorf("colour.White.String() = \"%s\"; want \"&H00FFFFFF\"", got)
	}
	if got := colour.Black.String(); got != "&H00000000" {
		t.Errorf("colour.Black.String() = \"%s\"; want \"&H00000000\"", got)
	}
}

func TestFromString(t *testing.T) {
	if got, err := colour.FromString("&H00FFFFFF"); err != nil {
		t.Errorf("Could not parse valid white colour: %s", err)
	} else if got != colour.White {
		t.Errorf("colour.FromString(\"&H00FFFFFF\") != colour.White")
	}
	if got, err := colour.FromString("&H00000000"); err != nil {
		t.Errorf("Could not parse valid black colour: %s", err)
	} else if got != colour.Black {
		t.Errorf("colour.FromString(\"&H00000000\") != colour.Black")
	}
	if _, err := colour.FromString(""); err == nil {
		t.Errorf("Could parse empty colour")
	}
	if _, err := colour.FromString("&H"); err == nil {
		t.Errorf("Could parse invalid colour (too short) \"&H\"")
	}
	if _, err := colour.FromString("&H000000000"); err == nil {
		t.Errorf("Could parse invalid colour (too long) \"&H000000000\"")
	}
	if _, err := colour.FromString("&HGHIJGHIJ"); err == nil {
		t.Errorf("Could parse invalid colour (not hex) \"&HGHIJGHIJ\"")
	}
}
