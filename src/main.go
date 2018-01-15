package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
)

const (
	chrome = `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`
	path   = `C:\Users\tkido\Dropbox\Kami Data\ブログ\test.txt`
)

func main() {
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
}
