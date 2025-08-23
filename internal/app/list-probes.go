// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/louisroyer/km-probe/internal/probes"
)

type ListProbes struct {
}

func NewListProbes() *ListProbes {
	return &ListProbes{}
}

type prb struct {
	name string
	desc string
}

const (
	AnsiUnderline = "\033[4m"
	AnsiReset     = "\033[0m"
)

func (p prb) Println(namelen int, b *strings.Builder, underline bool) {
	if underline {
		b.WriteString(AnsiUnderline)
	}
	b.WriteString(p.name)
	if underline {
		b.WriteString(AnsiReset)
	}
	for _ = range namelen - len(p.name) {
		b.WriteString(" ")
	}
	b.WriteString("\t")
	if underline {
		b.WriteString(AnsiUnderline)
	}
	b.WriteString(p.desc)
	if underline {
		b.WriteString(AnsiReset)
	}
	fmt.Println(b.String())
	b.Reset()
}

func (l *ListProbes) Run(ctx context.Context) error {
	list := make([]prb, 0)
	header := prb{name: "Name", desc: "Description"}
	namelen := len(header.name)
	for _, pf := range probes.AvailableProbes() {
		item := prb{name: pf.Name(), desc: pf.Description()}
		list = append(list, item)
		if len(item.name) > namelen {
			namelen = len(item.name)
		}
	}
	b := strings.Builder{}
	header.Println(namelen, &b, true)

	for _, item := range list {
		item.Println(namelen, &b, false)
	}
	return nil
}
