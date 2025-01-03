// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
	"strings"
)

func (p *Probe) CheckAutomation(ctx context.Context) error {

	for _, line := range p.LyricsFile.Events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if strings.HasPrefix(line, "Comment: ") {
				p.Report.Pass("automation")
				return nil
			}
		}
	}
	p.Report.Fail("automation")
	return nil
}
