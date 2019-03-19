package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Event handles file events. Currently, when a CREATE operation is seen, the
// .png file is added to the output file. Otherwise, the event is ignored.
func Event(event fsnotify.Event, ok bool, output *Output) {
	if !ok || event.Op != fsnotify.Create {
		return
	}

	// verify that it's a png file
	if len(event.Name) < 4 || event.Name[len(event.Name)-4:] != ".png" {
		return
	}

	// read png data
	data, err := ioutil.ReadFile(event.Name)
	if err != nil {
		log.Panicln(err)
	}

	// add image to output
	if err = output.Add(data); err != nil {
		log.Panicln(err)
	}
}

// Error is called when there's an error with the file notification loop.
func Error(err error, ok bool) {
	if !ok || err == nil {
		return
	}

	log.Panicln(err)
}

// Handle is just an infinite loop that handle file events and errors. Events
// are sent to Event and errors are sent to Error.
//
// Output is the singular output file where the screenshots are stored.
func Handle(watcher *fsnotify.Watcher, output *Output) {
	for {
		select {
		case event, ok := <-watcher.Events:
			Event(event, ok, output)
		case err, ok := <-watcher.Errors:
			Error(err, ok)
		}
	}
}

// Cleanup is ran when the user Ctrl+C's out of the application. It primarily
// just makes sure the output file is saved to disk.
func Cleanup(stop chan os.Signal, output *Output) {
	// wait for sigterm and sigint
	<-stop

	output.Save()
	log.Println("output saved.")

	// exit the application
	os.Exit(0)
}

func init() {
	// verify that the command line arguments are present
	if len(os.Args) < 3 {
		log.Fatalln("./swipe [watch dir] [output file]")
	}
}

func main() {
	var (
		output  *Output = NewOutput(os.Args[2])
		watcher *fsnotify.Watcher
		err     error
		stop    chan os.Signal = make(chan os.Signal)
	)

	// register stop channel with relevant signals
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	// cleanup function for program exit
	go Cleanup(stop, output)

	// initialize watcher
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// add the directories to watch
	if err = watcher.Add(os.Args[1]); err != nil {
		log.Fatalln(err)
	}

	// handle the events and errors
	Handle(watcher, output)
}
