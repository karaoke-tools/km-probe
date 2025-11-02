// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package lyrics_test

import (
	"slices"
	"testing"

	"github.com/karaoke-tools/km-probe/internal/ass/lyrics"
)

// TestKfLen tests `Line.KfLen`
func TestKfLen(t *testing.T) {
	line := lyrics.Line{
		TagsSplit: []string{
			"{\\k1}", // not a kf
			"test1",
			"{\\kf2}",
			"test2",
			"{\\kf3}",
			" ", // empty syllable
			"{\\kf4}",
			" test4", // syllable starting with space
		},
	}
	want := []int{2, 4}
	result := line.KfLen()
	if !slices.Equal(result, want) {
		t.Errorf(`Line.KfLen, result: "%v", want: "%v"`, result, want)
	}
}
