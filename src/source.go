package main

import (
	"os"
	"path/filepath"
	"sort"
)

func getSource() string {
	if _, err := os.Stat(namiPath); err == nil {
		return namiPath
	}
	txts, _ := filepath.Glob(globPath)
	sort.Slice(txts, func(i, j int) bool {
		f1, _ := os.Stat(txts[i])
		f2, _ := os.Stat(txts[j])
		return f1.ModTime().After(f2.ModTime())
	})
	return txts[0]
}
