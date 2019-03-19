package main

import (
	"fmt"
	"os/user"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// start the watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	// get current user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// block the app until quit (go routine nonblocking)
	done := make(chan bool)

	// goroutine for watching the folder
	go func() {
		for {
			select {
			// handle file event
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if event.Op == fsnotify.Create {
					// file placed on the desktop
					fmt.Println(event)
				}
			// handle error event
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				fmt.Println("error:", err)
			}
		}
	}()

	// TODO: other operating systems
	//
	// watches the current user's desktop directory
	if err = watcher.Add("/Users/" + user.Username + "/Desktop"); err != nil {
		panic(err)
	}

	<-done
}
