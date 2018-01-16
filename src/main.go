package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
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
	source := txts[0]
	
	// convert()
	exec.Command(chrome, html).Run()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(source)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	log.Println("Watching...")
	updated := time.Now()
Loop:
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write != 0 {
				now := time.Now()
				fmt.Println(ev)
				// "Write" event within 1.0 second is regarded as duplicated.
				if now.Sub(updated) > time.Second {
					// convert()
					fmt.Println("Converted!!")
					updated = now
				}
			}
		case err = <-watcher.Errors:
			log.Fatal(err)
		case s := <- ch:
			if s == syscall.SIGINT {
				break Loop
			}
		}
	}
	fmt.Println("Loop End")
	exec.Command(editor, result).Run()
	
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
