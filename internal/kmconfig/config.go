// Copyright Louis Royer. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package kmconfig

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type KmConfig struct {
	System System `yaml:"System"`
}

type System struct {
	Repositories []Repository `yaml:"Repositories"`
}

type Repository struct {
	Name               string `yaml:"Name"`
	Online             bool   `yaml:"Online"`
	Enabled            bool   `yaml:"Enabled"`
	Secure             bool   `yaml:"Secure"`
	SendStats          bool   `yaml:"SendStats"`
	BaseDir            string `yaml:"BaseDir"`
	Update             bool   `yaml:"Update"`
	AutoMediaDownloads string `yaml:"AutoMediaDownloads"`
	MaintainerMode     bool   `yaml:"MaintainerMode"`
	Path               Path   `yaml:"Path"`
}

type Path struct {
	Medias []string `yaml:"Medias"`
}

func ParseConf(file string) (*KmConfig, error) {
	var conf KmConfig
	path, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
