// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

const (
	BuildCommit = "dev"
)

// File is the configuration object.
var File *ini.File

func Init() error {
	configFile := os.Getenv("NEKOPIXEL_CONFIG_PATH")
	if configFile == "" {
		configFile = "conf/app.ini"
	}

	var err error
	File, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, configFile)
	if err != nil {
		return errors.Wrapf(err, "parse %q", configFile)
	}

	if err := File.Section("app").MapTo(&App); err != nil {
		return errors.Wrap(err, "map 'app'")
	}

	return nil
}
