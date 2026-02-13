// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package printer

import (
	"context"
	"sync"

	"github.com/karaoke-tools/km-probe/internal/lints"
)

// Printer outputs aggretator reports to stdout
type Printer interface {
	Encode(ctx context.Context, a *lints.Aggregator) error
	Aggregator() *lints.Aggregator
}

type BasePrinter struct {
	ready          chan struct{}
	aggregatorPool sync.Pool // allows reusing memory when creating aggregator
}

func NewBasePrinter() *BasePrinter {
	ready := make(chan struct{}, 1) // single simultaneous print because we print to stdout
	ready <- struct{}{}
	return &BasePrinter{
		ready: ready,
		aggregatorPool: sync.Pool{
			New: func() any {
				return lints.NewAggregator()
			},
		},
	}
}

func (p *BasePrinter) setReady() {
	p.ready <- struct{}{}
}

func (p *BasePrinter) Aggregator() *lints.Aggregator {
	return p.aggregatorPool.Get().(*lints.Aggregator)
}
