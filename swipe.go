package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// TODO: use dependency injection
var O = NewOutput("out.html")

type Output struct {
	Filename string
	Buffer   bytes.Buffer
}

func NewOutput(filename string) *Output {
	o := &Output{
		Filename: filename,
	}

	o.Buffer.WriteString(HTML_HEAD)

	return o
}

func (o *Output) Save() {
	o.Buffer.WriteString(HTML_TAIL)

	ioutil.WriteFile(o.Filename, o.Buffer.Bytes(), 0644)
}

func (o *Output) Add(b []byte) error {
	o.Buffer.WriteString(HTML_IMG_HEAD)
	o.Buffer.WriteString(base64.StdEncoding.EncodeToString(b))
	o.Buffer.WriteString(HTML_IMG_TAIL)

	return nil
}

func Event(event fsnotify.Event, ok bool) {
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
	if err = O.Add(data); err != nil {
		log.Panicln(err)
	}
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
	O.Save()

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
