// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package probe

import (
	"fmt"
)

type Report struct {
	DisplayName string
	Content     map[string]bool
}

func NewReport(kara_displayname string) *Report {
	return &Report{
		DisplayName: kara_displayname,
		Content:     make(map[string]bool),
	}
}

func (r *Report) Pass(name string) {
	r.Content[name] = true
}

func (r *Report) Fail(name string) {
	r.Content[name] = false
}

func (r *Report) String() string {
	return fmt.Sprintf("name: %s\n- automation: %t\n- resolution: %t\n- style: %t\n- live download probably allowed: %t\n- first contribution: %t\n",
		r.DisplayName,
		r.Content["automation"],
		r.Content["resolution"],
		r.Content["style"],
		r.Content[checkLiveDownloadProbablyAllowedKey],
		r.SuitableFirstContribution(),
	)
}

func (r *Report) SuitableFirstContribution() bool {
	if r.Content["automation"] && r.Content[checkLiveDownloadProbablyAllowedKey] && (!r.Content["resolution"] || !r.Content["style"]) {
		return true
	}
	return false
}
