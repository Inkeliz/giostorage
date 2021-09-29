//+build !js

package giostorage

import (
	"gioui.org/app"
	"io"
	"os"
	"path/filepath"
)

func (s *plainStorage) open(name string) (io.ReadCloser, error) {
	dir, err := app.DataDir()
	if err != nil {
		return nil, err
	}
	return os.Open(filepath.Join(dir, s.name, name))
}

func (s *plainStorage) create(name string) (io.WriteCloser, error) {
	dir, err := app.DataDir()
	if err != nil {
		return nil, err
	}
	dir = filepath.Join(dir, s.name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(filepath.Join(dir, name))
}
