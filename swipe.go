package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func Event(event fsnotify.Event, ok bool) {
	if !ok || event.Op != fsnotify.Create {
		return
	}

	log.Println(event)
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
