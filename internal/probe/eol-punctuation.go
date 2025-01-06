// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/lyrics"
)

const checkNoEolPunctuationKey = "eol-punctuation"

func (p *Probe) checkNoEolPunctuation(ctx context.Context) error {
	for _, line := range p.Lyrics.Events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if (line.Type != lyrics.Format) && (!strings.HasPrefix(line.Effect, "template")) {
				l := line.Text.StripTags()
				if strings.HasSuffix(l, ".") || strings.HasSuffix(l, ",") {
					if strings.HasSuffix(l, "...") {
						continue
					}
					p.Report.Fail(checkNoEolPunctuationKey)
					return nil
				}
			}
		}
	}
	p.Report.Pass(checkNoEolPunctuationKey)
	return nil
}
