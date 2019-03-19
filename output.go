package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

// Output holds the current output html file and associated methods.
type Output struct {
	Filename string
	Buffer   bytes.Buffer
}

// NewOutput sets the filename and returns an output pointer.
func NewOutput(filename string) *Output {
	o := &Output{
		Filename: filename,
	}

	o.Buffer.WriteString(HTML_HEAD)

	return o
}

// Save writes the output file to disk using the filename set in NewOutput().
func (o *Output) Save() {
	o.Buffer.WriteString(HTML_TAIL)

	ioutil.WriteFile(o.Filename, o.Buffer.Bytes(), 0644)
}

// Add adds the PNG byte data as a base64 image to the html file.
func (o *Output) Add(b []byte) error {
	o.Buffer.WriteString(HTML_IMG_HEAD)
	o.Buffer.WriteString(base64.StdEncoding.EncodeToString(b))
	o.Buffer.WriteString(HTML_IMG_TAIL)

	return nil
}
