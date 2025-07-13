// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"
)

var checkScaledBorderAndShadowKey = "scaled-border-and-shadow"

func (p *Probe) checkScaledBorderAndShadow(ctx context.Context) error {
	if p.Lyrics.ScriptInfo.ScaledBorderAndShadow {
		p.Report.Pass(checkScaledBorderAndShadowKey)
		return nil
	}
	p.Report.Fail(checkScaledBorderAndShadowKey)
	return nil
}
