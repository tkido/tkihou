package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	txts, err := filepath.Glob(path)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(txts, func(i, j int) bool {
		f1, _ := os.Stat(txts[i])
		f2, _ := os.Stat(txts[j])
		return f1.ModTime().After(f2.ModTime())
    })
	fmt.Println(txts[0])
	
	/*
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		// log.Print(strconv.Quote(s.Text()))
	}
	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}
	exec.Command(chrome, path).Run()
	exec.Command(editor, path).Run()
	*/
}
