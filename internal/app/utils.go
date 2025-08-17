// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"os"
	"path/filepath"

	"github.com/louisroyer/km-probe/internal/kmconfig"
)

var subpaths = []string{"karaokemugen-app", "karaokemugen-server", "km-server"}

func kmDataPath() (string, error) {
	kmDataPaths, err := xdgDataPaths()
	if err != nil {
		return "", err
	}
	// fallback on `..`
	// yes, this is very weird, but I don't judge you
	// it will make sense if you run it in
	// the exact same way km-server documentation says
	// (i.e. from the kmserver git repository)
	kmDataPaths = append(kmDataPaths, "..")
	for _, s := range subpaths {
		for _, p := range kmDataPaths {
			path := filepath.Join(p, s, "app/")
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}
	return "", ErrDataRepositoryNotFound
}

func kmConfigPath() (string, error) {
	kmConfPaths := make([]string, 0, 6)
	// 1. load config from "historical" KM config-file path (`XDG_DATA_HOME`/`XDG_DATA_DIRS`)
	xdgData, err := xdgDataPaths()
	if err != nil {
		return "", err
	}
	kmConfPaths = append(kmConfPaths, xdgData...)
	// 2. maybe we are in the future, and now KM is "really" following
	// XDG spec for the config-file path (`XDG_CONFIG_HOME`/`XDG_CONFIG_DIRS`)?
	xdgConfig, err := xdgConfigPaths()
	if err != nil {
		return "", err
	}
	kmConfPaths = append(kmConfPaths, xdgConfig...)

	// fallback on `..`
	// yes, this is very weird, but I don't judge you
	// it will make sense if you run it in
	// the exact same way km-server documentation says
	// (i.e. from the kmserver git repository)
	kmConfPaths = append(kmConfPaths, "..")
	for _, s := range subpaths {
		for _, p := range kmConfPaths {
			path := filepath.Join(p, s, "app/")
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}
	return "", ErrConfigNotFound
}

func loadConf() (*kmconfig.KmConfig, error) {
	path, err := kmConfigPath()
	if err != nil {
		return nil, err
	}
	if kmConfig, err := kmconfig.ParseConf(filepath.Join(path, "config.yml")); err == nil {
		return kmConfig, nil
	}
	return nil, ErrConfigNotFound
}
