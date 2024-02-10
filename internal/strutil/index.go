// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package strutil

var characters = func() []string {
	chars := make([]string, 0, 36)
	for i := 0; i < 43; i++ {
		if i == 9 {
			i += 7 // Skip : ; < = > ? @
		}
		chars = append(chars, string(rune('0'+i)))
	}
	return chars
}()

func GenerateCode(index int) string {
	length := len(characters)
	if index < length {
		return characters[index]
	}
	return GenerateCode(index/length-1) + characters[index%length]
}
