// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	ink "github.com/CatInBeard/inkview"
)

const defaultFontSize = 20

func main() {
	defer func() {
		if r := recover(); r != nil {
			var pcs [3]uintptr
			n := runtime.Callers(3, pcs[:])
			if n > 0 {
				frames := runtime.CallersFrames(pcs[:n])
				var stackInfo []string
				for frame, more := frames.Next(); more; frame, more = frames.Next() {
					stackInfo = append(stackInfo, fmt.Sprintf("%s:%d", frame.File, frame.Line))
				}
				ink.Errorf("Err", fmt.Sprintf("%v\n%s", r, strings.Join(stackInfo, "\n")))
			} else {
				ink.Errorf("Err", fmt.Sprintf("%v", r))
			}
		}
	}()

	app := &DispatcherApp{fontH: defaultFontSize}

	if err := ink.Run(app); err != nil {
		log.Fatal(err)
	}

}
