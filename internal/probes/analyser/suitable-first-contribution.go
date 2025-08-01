// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package analyser

import (
	"context"

	"github.com/louisroyer/km-probe/internal/probes/report"
)

type SuitableFirstContribution struct {
	baseAnalyser
}

func NewSuitableFirstContribution(r map[string]report.Report) Analyser {
	return &SuitableFirstContribution{
		newAnalyser("suitable-first-contribution", r),
	}
}
func (a *SuitableFirstContribution) Run(ctx context.Context) (report.Report, error) {
	critical := []string{
		"live-download",
		"resolution",
		"automation",
	}
	for _, c := range critical {
		if r, ok := a.reports[c]; !ok {
			return report.Skip(), nil
		} else if !r.Result() {
			return report.Info(false), nil
		}
	}

	scoring := [][]string{
		// style issues
		[]string{"style-single-white", "style-black-border"}, // minor issues
		[]string{"resolution"},                               // can imply re-splitting some parts
		// lyrics issues
		[]string{"double-consonnant"},
	}

	badness := 0
	for _, s := range scoring {
		local_badness := 0
		for _, sb := range s {
			if r, ok := a.reports[sb]; !ok {
				return report.Skip(), nil
			} else if !r.Result() {
				local_badness++
			}
		}
		if local_badness > 0 {
			badness++
		}
	}

	return report.Info(badness == 1), nil
}
