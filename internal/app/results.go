// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"fmt"

	"github.com/louisroyer/km-probe/internal/probe"
)

type Results struct {
	Path       string
	Automation *probe.Report
	Resolution *probe.Report
}

func NewResults(path string) *Results {
	return &Results{Path: path}
}

func (r *Results) String() string {
	return fmt.Sprintf("name: %s\n- automation: %s\n- resolution: %s\n- first contribution: %s\n",
		r.Path,
		r.Automation.Content["pass"],
		r.Resolution.Content["pass"],
		r.SuitableFirstContribution(),
	)
}

func (r *Results) SuitableFirstContribution() string {
	//FIXME: return type
	if r.Automation.Content["pass"] == "true" && r.Resolution.Content["pass"] == "false" {
		return "true"
	}
	return "false"
}
