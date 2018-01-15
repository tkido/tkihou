package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/fsnotify/fsnotify"
)

func getSource() string {
	txts, err := filepath.Glob(path)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(txts, func(i, j int) bool {
		f1, _ := os.Stat(txts[i])
		f2, _ := os.Stat(txts[j])
		return f1.ModTime().After(f2.ModTime())
    })
	return txts[0]
}

func main() {
	source := getSource()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	defer exec.Command(editor, source).Run()

	err = watcher.Add(source)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching...")
	updated := time.Now()
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write != 0 {
				now := time.Now()
				fmt.Println(ev)
				// "Write" event within 0.5 second to the same file is regarded as duplicated.
				if now.Sub(updated) > time.Second/2 {
					// convert
					fmt.Println("Converted!!")
					updated = now
				}
			}
		case err = <-watcher.Errors:
			log.Fatal(err)
		}
	}

	
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
