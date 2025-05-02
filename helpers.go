// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import "unicode/utf8"

func firstNRunes(s string, n int) string {
	runes := []rune(s)
	if n > len(runes) {
		return s
	}
	if n <= 0 {
		return s

	}
	return string(runes[:n])
}

func countRealChar(str string) int {
	return utf8.RuneCount([]byte(str))
}
