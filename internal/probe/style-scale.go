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

func (p *Probe) checkStyleScale(ctx context.Context) error {
	for _, line := range p.Lyrics.Styles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				s, err := style.Parse(strings.TrimPrefix(line, "Style: "))
				if err != nil {
					return err
				}
				if (s.ScaleX != "100") || (s.ScaleY != "100") {
					p.Report.Fail("style-scale")
					return nil
				}
			}
		}
	}
	p.Report.Pass("style-scale")
	return nil
}
