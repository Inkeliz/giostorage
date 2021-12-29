//go:build js
// +build js

package giostorage

import "encoding/json"

var (
	Encoder = json.NewEncoder
	Decoder = json.NewDecoder
)
