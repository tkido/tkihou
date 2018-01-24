package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	source := getSource()
	convert(source)
	exec.Command(chrome, rstHTML).Run()

	if flags.Watch {
		watch(source)
	}
	exec.Command(editor, rstTxt).Start()
}

func watch(source string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	watcher.Add(source)
	updated := time.Now()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	quit := make(chan bool)
	go func() {
		fmt.Scanln()
		quit <- true
	}()

	log.Println("Watching...")
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write != 0 {
				time.Sleep(time.Second / 10) // wait 0.1s for workaround
				now := time.Now()
				// "Write" event within 0.5 second is regarded as duplicated.
				if now.Sub(updated) > time.Second/2 {
					convert(source)
					updated = now
					log.Println("Converted!")
				}
			}
		case err := <-watcher.Errors:
			log.Fatal(err)
		case s := <-ch:
			if s == syscall.SIGINT {
				return
			}
		case <-quit:
			return
		}
	}
}
