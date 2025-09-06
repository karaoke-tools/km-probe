// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package printer

import (
	"context"
	"encoding/json"
	"os"

	"github.com/karaoke-tools/km-probe/internal/probes"
)

type JsonPrinter struct {
	*BasePrinter
	e *json.Encoder
}

func NewJsonPrinter() Printer {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return &JsonPrinter{
		BasePrinter: NewBasePrinter(),
		e:           encoder,
	}
}

func (p *JsonPrinter) Encode(ctx context.Context, a *probes.Aggregator) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ready:
		defer p.setReady()
		defer p.aggregatorPool.Put(a)
		if err := p.e.Encode(a); err != nil {
			return err
		}
	}
	return nil

}
