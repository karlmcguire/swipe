package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	// "github.com/jung-kurt/gofpdf"
)

// start:
// 	- scan the directory, add each new file created to pdf
// stop:
// 	- finish generating the pdf, print out file location

func Event(event fsnotify.Event, ok bool) {
	if !ok || event.Op != fsnotify.Create {
		return
	}

	log.Println(event.Name[len(event.Name)-4:])
}

func Error(err error, ok bool) {
	if !ok || err == nil {
		return
	}

	log.Panicln(err)
}

func Handle(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			Event(event, ok)
		case err, ok := <-watcher.Errors:
			Error(err, ok)
		}
	}
}

func Cleanup(stop chan os.Signal) {
	// wait for sigterm and sigint
	<-stop

	log.Println("stopping")

	// exit the application
	os.Exit(0)
}

func init() {
	// channel for sigterm and sigint
	stop := make(chan os.Signal)

	// register channel
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	// cleanup function
	go Cleanup(stop)
}

func main() {
	var (
		watcher *fsnotify.Watcher
		err     error
	)

	// initialize watcher
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// add the directories to watch
	if err = watcher.Add(PATH_WATCH); err != nil {
		log.Fatalln(err)
	}

	// handle the events and errors
	Handle(watcher)
}
