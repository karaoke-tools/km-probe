// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"context"

	"github.com/louisroyer/km-probe/internal/karadata"
	"github.com/louisroyer/km-probe/internal/probes/report"
)

type Probe interface {
	Name() string
	PreRun(ctx context.Context, KaraData *karadata.KaraData) (bool, string, error)
	Run(ctx context.Context, KaraData *karadata.KaraData) (report.Report, error)
	Description() string
}
