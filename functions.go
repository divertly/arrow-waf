package main

import (
	"io"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
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

func LoadRuleList(fname string) (*CRSRuleList, error) {
	rl := &CRSRuleList{}
	fd, err := os.OpenFile(fname, os.O_RDONLY, 0600) // #nosec:G304
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &rl)
	if err != nil {
		return nil, err
	}
	return rl, nil
}
