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
			"{\\kf5}",
			"test5|test5", // syllable with furigana
			"{\\kf6}",
			"test6", // continuation of previous kf
			"{\\kf7}",
			"test7", // this should not include 7 because (5+6+7)/3 == 5 and |7-5| > 1
			"{\\kf8}",
			"|test8",
			"{\\k404}", // this should stop now
			"test404",
			"{\\kf9}",
			"test9", // this should not be grouped with previous 8

		},
	}
	want := []int{2, 4, 11, 7, 8, 9}
	result := line.KfLen()
	if !slices.Equal(result, want) {
		t.Errorf(`Line.KfLen, result: "%v", want: "%v"`, result, want)
	}
}
