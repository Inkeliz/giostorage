package giostorage

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

func TestNewEncryptedStorage(t *testing.T) {
	key := make([]byte, EncryptedKeySize)
	rand.Read(key)

	storage, err := NewEncryptedStorage("name", key)
	if err != nil {
		t.Error(err)
		return
	}

	fileCreate, err := storage.Create("test")
	if err != nil {
		t.Error(err)
		return
	}

	randData := make([]byte, 16000)
	rand.Read(randData)

	n, err := io.Copy(fileCreate,  bytes.NewReader(randData))
	if err != nil {
		t.Error(err)
		return
	}

	if int(n) != len(randData) {
		t.Error("invalid len")
		return
	}

	if err := fileCreate.Close(); err != nil {
		t.Error(err)
		return
	}

	fileOpen, err := storage.Open("test")
	if err != nil {
		t.Error(err)
		return
	}

	allData, err := io.ReadAll(fileOpen)
	if err != nil {
		t.Error(err)
		return
	}

	if err := fileOpen.Close(); err != nil {
		t.Error(err)
		return
	}

	plainStorage := NewPlainStorage("name")

	filePlain, err := plainStorage.Open("test")
	if err != nil {
		t.Error(err)
		return
	}

	plainData, err := io.ReadAll(filePlain)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Equal(plainData, randData) {
		t.Error("not encrypted")
		return
	}

	if !bytes.Equal(allData, randData) {
		t.Error("invalid data")
		return
	}
}