// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ass

type ScriptInfo struct {
	PlayResX              uint32
	PlayResY              uint32
	ScaledBorderAndShadow bool
}

func (s *ScriptInfo) SetScaledBorderAndShadow(v string) error {
	switch v {
	case "yes":
		s.ScaledBorderAndShadow = true
	case "no":
		s.ScaledBorderAndShadow = false
	default:
		return ErrMalformedFile
	}
	return nil
}
