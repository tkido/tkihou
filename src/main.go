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
	convert(getSource())
	exec.Command(chrome, rstHTML).Run()

	if flags.Watch {
		watch()
	}
	exec.Command(editor, rstTxt).Start()
}

func watch() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	watcher.Add(watchPath)
	updated := time.Now()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	quit := make(chan bool)
	go func() {
		fmt.Scanln()
		quit <- true
	}()

	log.Println("Watching...")
	// count := 0
	for {
		select {
		case ev := <-watcher.Events:
			//count++
			//log.Printf("%d 回目のイベント\n", count)
			time.Sleep(time.Second / 10) // wait 0.1s for workaround
			if ev.Op&fsnotify.Chmod != 0 {
				log.Println("Chmod")
			} else {
				/*
					if ev.Op&fsnotify.Write != 0 {
						log.Println("Write")
					} else if ev.Op&fsnotify.Create != 0 {
						log.Println("Create")
					} else if ev.Op&fsnotify.Remove != 0 {
						log.Println("Remove")
					} else if ev.Op&fsnotify.Rename != 0 {
						log.Println("Rename")
					}
				*/
				now := time.Now()
				// event within 0.5 second is regarded as duplicated.
				if now.Sub(updated) > time.Second/2 {
					convert(getSource())
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
