package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

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
