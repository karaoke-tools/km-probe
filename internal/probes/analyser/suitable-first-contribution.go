// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package analyser

import (
	"context"

	"github.com/louisroyer/km-probe/internal/probes/report"
	"github.com/louisroyer/km-probe/internal/probes/report/result"
)

type SuitableFirstContribution struct {
	baseAnalyser
}

func NewSuitableFirstContribution(r map[string]report.Report) Analyser {
	return &SuitableFirstContribution{
		newAnalyser("kara-moe.suitable-first-contribution", r),
	}
}
func (a *SuitableFirstContribution) Run(ctx context.Context) (report.Report, error) {
	critical := []string{
		"kara-moe.live-download",
		"system.resolution",
		"system.automation",
	}
	for _, c := range critical {
		if r, ok := a.reports[c]; !ok {
			return report.Skip("a report is missing"), nil
		} else if r.Result() != result.Passed {
			return report.Info(false), nil
		}
	}

	scoring := [][]string{
		// style issues
		[]string{"kara-moe.style-single-white", "system.style-black-border"}, // minor issues
		[]string{"system.resolution"},                                        // can imply re-splitting some parts
		// lyrics issues
		[]string{"kara-moe.double-consonant"},
	}

	badness := 0
	for _, s := range scoring {
		local_badness := 0
		for _, sb := range s {
			if r, ok := a.reports[sb]; !ok {
				return report.Skip("missing report"), nil
			} else if r.Result() != result.Failed {
				local_badness++
			}
		}
		if local_badness > 0 {
			badness++
		}
	}

	return report.Info(badness == 1), nil
}
