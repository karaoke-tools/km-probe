// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package baseprobe

type EnabledByDefault struct{}

func (p EnabledByDefault) Enabled() bool {
	return true
}
