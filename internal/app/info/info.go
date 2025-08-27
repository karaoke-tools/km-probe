// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package info

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/louisroyer/km-probe/internal/app/setup"
	"github.com/louisroyer/km-probe/internal/probes"

	"github.com/urfave/cli/v2"
)

type ListProbes struct {
	*setup.Setup
}

func FromCli(ctx *cli.Context) *ListProbes {
	return &ListProbes{Setup: setup.FromCli(ctx)}
}

type prb struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

const (
	AnsiUnderline = "\033[4m"
	AnsiReset     = "\033[0m"
)

func (p prb) Println(namelen int, b *strings.Builder, underline bool) {
	if underline {
		b.WriteString(AnsiUnderline)
	}
	b.WriteString(p.Name)
	if underline {
		b.WriteString(AnsiReset)
	}
	for _ = range namelen - len(p.Name) {
		b.WriteString(" ")
	}
	b.WriteString("\t")
	if underline {
		b.WriteString(AnsiUnderline)
	}
	b.WriteString(p.Desc)
	if underline {
		b.WriteString(AnsiReset)
	}
	fmt.Println(b.String())
	b.Reset()
}

func (l *ListProbes) Run(ctx context.Context) error {
	if l.OutputJson {
		return l.RunJson(ctx)
	}
	return l.RunTxt(ctx)
}

func (l *ListProbes) RunTxt(ctx context.Context) error {
	list := make([]prb, 0)
	header := prb{Name: "Name", Desc: "Description"}
	namelen := len(header.Name)
	for _, pf := range probes.AvailableProbes() {
		item := prb{Name: pf.Name(), Desc: pf.Description()}
		list = append(list, item)
		if len(item.Name) > namelen {
			namelen = len(item.Name)
		}
	}
	b := strings.Builder{}
	header.Println(namelen, &b, l.Color)

	for _, item := range list {
		item.Println(namelen, &b, false)

	}
	return nil
}

func (l *ListProbes) RunJson(ctx context.Context) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	for _, pf := range probes.AvailableProbes() {
		encoder.Encode(prb{Name: pf.Name(), Desc: pf.Description()})
	}
	return nil
}

func RunFromCli(ctx *cli.Context) error {
	return FromCli(ctx).Run(ctx.Context)
}
