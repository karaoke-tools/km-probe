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

func (p prb) Println(namelen int, b *strings.Builder) {
	b.WriteString(p.name)
	for _ = range namelen - len(p.name) {
		b.WriteString(" ")
	}
	b.WriteString("\t")
	b.WriteString(p.desc)
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
	header.Println(namelen, &b)

	for _, item := range list {
		item.Println(namelen, &b)
	}
	return nil
}
