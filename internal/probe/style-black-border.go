// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
)

func (p *Probe) checkStyleBlackBorder(ctx context.Context) error {
	for _, line := range p.Lyrics.Styles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") { // we don't care about furigana styles for the now
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return err
				}
				if s.OutlineColour != "&H00000000" {
					// secondary color must be black
					p.Report.Fail("style")
					return nil
				}
				break
			}
		}
	}
	p.Report.Pass("style-black-border")
	return nil
}
