// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package analyser

import (
	"context"

	"github.com/louisroyer/km-probe/internal/probes/report"
)

type NewAnalyserFunc func(map[string]report.Report) Analyser

type baseAnalyser struct {
	name    string
	reports map[string]report.Report
}

func newAnalyser(name string, reports map[string]report.Report) baseAnalyser {
	return baseAnalyser{
		name:    name,
		reports: reports,
	}
}

func (p *baseAnalyser) Name() string {
	return p.name
}

type Analyser interface {
	Name() string
	Run(ctx context.Context) (report.Report, error)
}
