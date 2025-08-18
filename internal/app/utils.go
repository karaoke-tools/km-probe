// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/louisroyer/km-probe/internal/kmconfig"

	"github.com/adrg/xdg"
)

var subpaths = []string{"karaokemugen-app", "karaokemugen-server", "km-server"}

const configpath = "app/config.yml"

func searchKmDataDirPath(baseDir string) (string, error) {
	if filepath.IsAbs(baseDir) {
		if _, err := os.Stat(baseDir); errors.Is(err, fs.ErrNotExist) {
			return baseDir, ErrDataRepositoryNotFound
		}
		return baseDir, nil
	}

	for _, s := range subpaths {
		path, err := xdg.SearchDataFile(filepath.Join(s, "app", baseDir))
		if err == nil {
			return path, nil
		}
	}
	// fallback (if this is run from a kmserver git repository)
	path := filepath.Join("app", baseDir)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", ErrDataRepositoryNotFound
}

func searchKmConfigFilePath() (string, error) {
	// Search on XDG basedir compliant paths for config
	// (just in case KM start to be compliant)
	for _, s := range subpaths {
		path, err := xdg.SearchConfigFile(filepath.Join(s, configpath))
		if err == nil {
			return path, nil
		}
	}
	// Search paths really used by KM
	for _, s := range subpaths {
		path, err := xdg.SearchDataFile(filepath.Join(s, configpath))
		if err == nil {
			return path, nil
		}
	}
	// fallback (if this is run from a kmserver git repository)
	path := configpath
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", ErrConfigNotFound
}

func loadConf() (*kmconfig.KmConfig, error) {
	path, err := searchKmConfigFilePath()
	if err != nil {
		return nil, err
	}
	if kmConfig, err := kmconfig.ParseConf(path); err == nil {
		return kmConfig, nil
	}
	return nil, ErrConfigNotFound
}
