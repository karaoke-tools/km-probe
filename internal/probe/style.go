// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"

	"github.com/louisroyer/km-probe/internal/ass/style"
	"github.com/louisroyer/km-probe/internal/karajson"
)

func (p *Probe) checkStyle(ctx context.Context) error {
	for _, tag := range p.KaraJson.Data.Tags.Misc {
		if tag == karajson.GroupSinging {
			p.Report.Pass("style")
			return nil
		}
	}
	nb_styles := 0
	for _, line := range p.Lyrics.Styles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if strings.HasPrefix(line, "Style: ") && !strings.Contains(line, "-furigana") {
				nb_styles += 1
				if nb_styles > 1 {
					// for the moment, we focus on single style karaoke
					p.Report.Fail("style")
					return nil
				}
			}
		}
	}
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
				if s.SecondaryColour == "&H00FFFFFF" {
					p.Report.Pass("style")
					return nil
				}
				break
			}
		}
	}
	p.Report.Fail("style")
	return nil
}
