// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func xdgDataPaths() ([]string, error) {
	xdgPaths := make([]string, 0, 3)

	xdgDataHome, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		usr, err := user.Current()
		if err != nil {
			return []string{}, err
		}
		xdgDataHome = filepath.Join(usr.HomeDir, ".local/share")
	}
	xdgPaths = append(xdgPaths, xdgDataHome)

	xdgDataDirs, ok := os.LookupEnv("XDG_DATA_DIRS")
	if !ok {
		xdgDataDirs = "/usr/local/share/:/usr/share"
	}
	xdgPaths = append(xdgPaths, strings.Split(xdgDataDirs, ":")...)
	return xdgPaths, nil
}

func xdgConfigPaths() ([]string, error) {
	xdgPaths := make([]string, 0, 2)

	xdgConfigHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		usr, err := user.Current()
		if err != nil {
			return []string{}, err
		}
		xdgConfigHome = filepath.Join(usr.HomeDir, ".config")
	}
	xdgPaths = append(xdgPaths, xdgConfigHome)

	xdgConfigDirs, ok := os.LookupEnv("XDG_CONFIG_DIRS")
	if !ok {
		xdgConfigDirs = "/etc/xdg"
	}
	xdgPaths = append(xdgPaths, strings.Split(xdgConfigDirs, ":")...)
	return xdgPaths, nil
}
