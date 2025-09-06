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

	"github.com/karaoke-tools/km-probe/internal/app/ansi"
	"github.com/karaoke-tools/km-probe/internal/app/setup"
	"github.com/karaoke-tools/km-probe/internal/probes"

	"github.com/urfave/cli/v2"
)

type ListProbes struct {
	*setup.Setup
}

func FromCli(ctx *cli.Context) *ListProbes {
	return &ListProbes{Setup: setup.FromCli(ctx)}
}

type prb struct {
	Name          string `json:"name"`
	Desc          string `json:"description"`
	Enabled       bool   `json:"enabled"`
	EnabledString string `json:"-"`
}

func enabledString(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

type color int

const (
	green color = iota
	red
	noColor
)

func (p prb) Println(namelen int, enabledlen int, b *strings.Builder, underline bool, color color) {
	switch color {
	case green:
		b.WriteString(ansi.Green)
	case red:
		b.WriteString(ansi.Red)
	}
	// Enabled
	if underline {
		b.WriteString(ansi.Underline)
	}
	b.WriteString(p.EnabledString)
	if underline {
		b.WriteString(ansi.Reset)
	}
	for _ = range enabledlen - len(p.EnabledString) {
		b.WriteString(" ")
	}

	b.WriteString("\t")

	// Name
	if underline {
		b.WriteString(ansi.Underline)
	}
	b.WriteString(p.Name)
	if underline {
		b.WriteString(ansi.Reset)
	}
	for _ = range namelen - len(p.Name) {
		b.WriteString(" ")
	}

	b.WriteString("\t")

	// Description
	if underline {
		b.WriteString(ansi.Underline)
	}
	b.WriteString(p.Desc)
	if underline || color != noColor {
		b.WriteString(ansi.Reset)
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
	header := prb{Name: "Name", Desc: "Description", EnabledString: "Status"}
	namelen, enabledlen := len(header.Name), len(header.EnabledString)
	for _, pf := range probes.AvailableProbes() {
		item := prb{Name: pf.Name(), Desc: pf.Description(), Enabled: pf.Enabled(), EnabledString: enabledString(pf.Enabled())}
		list = append(list, item)
		if len(item.Name) > namelen {
			namelen = len(item.Name)
		}
		if len(item.EnabledString) > enabledlen {
			enabledlen = len(item.EnabledString)
		}
	}
	b := strings.Builder{}
	header.Println(namelen, enabledlen, &b, l.Color, noColor)

	for _, item := range list {
		c := noColor
		if l.Color {
			if item.Enabled {
				c = green
			} else {
				c = red
			}
		}
		item.Println(namelen, enabledlen, &b, false, c)

	}
	return nil
}

func (l *ListProbes) RunJson(ctx context.Context) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	for _, pf := range probes.AvailableProbes() {
		encoder.Encode(prb{Name: pf.Name(), Desc: pf.Description(), Enabled: pf.Enabled()})
	}
	return nil
}

func RunFromCli(ctx *cli.Context) error {
	return FromCli(ctx).Run(ctx.Context)
}
