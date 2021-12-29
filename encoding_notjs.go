//go:build !js
// +build !js

package giostorage

import "encoding/gob"

var (
	Encoder = gob.NewEncoder
	Decoder = gob.NewDecoder
)
