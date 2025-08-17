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

type NewProbeFunc func(*karadata.KaraData) Probe

type Probe interface {
	Name() string
	Run(ctx context.Context) (report.Report, error)
	Description() string
}
