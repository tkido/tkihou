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
		log.Println("Watching...")
		watch(source)
	}
	exec.Command(editor, rstTxt).Run()
}

func watch(source string) {
	exec.Command(editor, source).Run()

	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	watcher.Add(source)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	updated := time.Now()

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write != 0 {
				now := time.Now()
				fmt.Println(ev)
				// "Write" event within 0.5 second is regarded as duplicated.
				if now.Sub(updated) > time.Second/2 {
					convert(source)
					fmt.Println("Converted!!")
					updated = now
				}
			}
		case err := <-watcher.Errors:
			log.Fatal(err)
		case s := <-ch:
			if s == syscall.SIGINT {
				return
			}
		}
	}
}
