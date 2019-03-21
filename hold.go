package main

import (
	"bytes"
	"os"
)

type Hold struct {
	Buffer bytes.Buffer
}

func (h *Hold) Write(d string) {
	h.Buffer.WriteString(d)
}

func (h *Hold) Store(n string) error {
	f, err := os.OpenFile(n, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = f.Write(h.Buffer.Bytes()); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	h.Buffer.Reset()

	return nil
}
