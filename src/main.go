package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func convert(src string) {
	f, err := os.Open(src)
	defer f.Close()
	if err != nil { log.Fatal(err) }
	s := bufio.NewScanner(f)
	if !s.Scan() { log.Fatal("No Title!!") }
	title := s.Text()
	buf := bytes.NewBufferString(fmt.Sprintf("<!--\n%s\n-->\n", title))
	for s.Scan() {
		buf.WriteString(s.Text())
	}
	if s.Err() != nil { log.Fatal(s.Err()) }
	content := buf.String()
	execute(title, content)
}

func main() {
	source := getSource()
	convert(source)
	exec.Command(chrome, rstHTML).Run()

	watcher, err := fsnotify.NewWatcher()
	if err != nil { log.Fatal(err) }
	defer watcher.Close()

	err = watcher.Add(source)
	if err != nil { log.Fatal(err) }

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
					convert(source)
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
	exec.Command(editor, rstTxt).Run()
}
