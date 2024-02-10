// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

var App struct {
	MaxHeight int      `ini:"max_height"`
	MaxWidth  int      `ini:"max_width"`
	Colors    []string `ini:"colors"`
}
