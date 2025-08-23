// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"encoding/json"
	"os"

	"github.com/louisroyer/km-probe/internal/probes"
)

type Printer struct {
	e     *json.Encoder
	ready chan struct{}
}

func NewPrinter() *Printer {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	ready := make(chan struct{}, 1)
	ready <- struct{}{}
	return &Printer{
		e:     encoder,
		ready: ready,
	}
}

func (p *Printer) setReady() {
	p.ready <- struct{}{}
}

func (p *Printer) Encode(ctx context.Context, a *probes.Aggregator) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ready:
		defer p.setReady()
		if err := p.e.Encode(a); err != nil {
			return err
		}
	}
	return nil

}
