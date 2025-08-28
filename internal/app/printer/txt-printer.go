// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package printer

import (
	"context"
	"fmt"
	"net/url"

	"github.com/louisroyer/km-probe/internal/app/ansi"
	"github.com/louisroyer/km-probe/internal/probes"
)

type TxtPrinter struct {
	*BasePrinter
	Hyperlink bool
	Color     bool
	BaseUri   string
}

func NewTxtPrinter(hyperlink bool, color bool, baseUri string) Printer {
	return &TxtPrinter{
		BasePrinter: NewBasePrinter(),
		Hyperlink:   hyperlink,
		Color:       color,
		BaseUri:     baseUri,
	}
}

func (p *TxtPrinter) Encode(ctx context.Context, a *probes.Aggregator) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ready:
		defer p.setReady()
		defer p.aggregatorPool.Put(a)
		if err := p.encodeAggregator(ctx, a); err != nil {
			return err
		}
	}
	return nil

}

func (p *TxtPrinter) encodeAggregator(ctx context.Context, a *probes.Aggregator) error {
	if p.Hyperlink {
		u, err := url.JoinPath(p.BaseUri, a.Kid.String())
		if err != nil {
			return err
		}
		fmt.Println(ansi.Link(u, a.Songname))
	} else {
		fmt.Println(a.Songname)
	}
	for k, r := range a.Reports {
		fmt.Printf("\t%s: %s\n", k, r.Result()) // TODO: alignment
		// TODO: color depending on severity
		// TODO: optional message on 2nd line with alignment
		// TODO: not display skipped/aborted: only the number
	}
	fmt.Println("") // empty line to separate aggregators
	return nil
}
