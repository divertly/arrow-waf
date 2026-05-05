package main

import (
	"runtime"
)

func this() (string, string) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "func", "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "func", "?"
	}

	return "func", fn.Name()
}
