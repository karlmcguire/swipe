package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Event handles file events. Currently, when a CREATE operation is seen, the
// .png file is added to the output file. Otherwise, the event is ignored.
func Event(eve fsnotify.Event, ok bool, h *Hold) {
	if !ok || eve.Op != fsnotify.Create {
		return
	}

	// verify that it's a png file
	if len(eve.Name) < 4 || eve.Name[len(eve.Name)-4:] != ".png" {
		return
	}

	// read png data
	data, err := ioutil.ReadFile(eve.Name)
	if err != nil {
		log.Panicln(err)
	}

	h.Write(HTML_IMG_HEAD)
	h.Write(base64.StdEncoding.EncodeToString(data))
	h.Write(HTML_IMG_TAIL)
	h.Store(os.Args[2])

	log.Println("saved '" + eve.Name + "' to disk.")
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
func Handle(w *fsnotify.Watcher, h *Hold) {
	for {
		select {
		case eve, ok := <-w.Events:
			Event(eve, ok, h)
		case err, ok := <-w.Errors:
			Error(err, ok)
		}
	}
}

// Cleanup is ran when the user Ctrl+C's out of the application. It primarily
// just makes sure the output file is saved to disk.
func Cleanup(s chan os.Signal, h *Hold) {
	<-s

	// save to disk
	h.Write(HTML_TAIL)
	h.Store(os.Args[2])
	log.Println("saved.")

	os.Exit(0)
}

func init() {
	if len(os.Args) < 3 {
		log.Fatalln("./swipe [watch dir] [output file]")
	}
}

func main() {
	h := new(Hold)
	h.Write(HTML_HEAD)

	{
		s := make(chan os.Signal)

		signal.Notify(s, syscall.SIGTERM)
		signal.Notify(s, syscall.SIGINT)

		go Cleanup(s, h)
	}

	{
		w, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalln(err)
		}
		defer w.Close()

		if err = w.Add(os.Args[1]); err != nil {
			log.Fatalln(err)
		}

		Handle(w, h)
	}
}
