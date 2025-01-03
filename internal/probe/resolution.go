// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
)

func (p *Probe) CheckResolution(ctx context.Context) error {
	if p.LyricsFile.ScriptInfo.PlayResX == 0 && p.LyricsFile.ScriptInfo.PlayResY == 0 {
		p.Report.Pass("resolution")
		return nil
	}
	p.Report.Fail("resolution")
	return nil
}
