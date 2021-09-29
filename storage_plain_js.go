package giostorage

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"syscall/js"
)

var (
	_LocalStorage = js.Global().Get("localStorage")
)

func (s *plainStorage) getName(name string) string {
	return s.name + "/" + name
}

func (s *plainStorage) open(name string) (io.ReadCloser, error) {
	v := _LocalStorage.Call("getItem", s.getName(name))
	if !v.Truthy() {
		return nil, os.ErrNotExist
	}

	r, err := hex.DecodeString(v.String())
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(r)), nil
}

func (s *plainStorage) create(name string) (io.WriteCloser, error) {
	return &file{name: s.getName(name), buffer: new(bytes.Buffer)}, nil
}

type file struct {
	name   string
	buffer *bytes.Buffer
}

func (f *file) Write(p []byte) (n int, err error) {
	return f.buffer.Write(p)
}

func (f *file) Close() error {
	_LocalStorage.Call("setItem", f.name, hex.EncodeToString(f.buffer.Bytes()))
	return nil
}
