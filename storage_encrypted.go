package giostorage

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/chacha20"
	"io"
)

type encryptedStorage struct {
	plain *plainStorage
	key   []byte
}

func (s *encryptedStorage) open(name string) (io.ReadCloser, error) {
	f, err := s.plain.open(name)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSizeX)
	if n, err := f.Read(nonce); n != chacha20.NonceSizeX || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("invalid nonce in file")
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(s.key, nonce)
	if err != nil {
		return nil, err
	}

	return &encryptRead{file: f, cipher: cipher}, nil
}

func (s *encryptedStorage) create(name string) (io.WriteCloser, error) {
	f, err := s.plain.create(name)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, chacha20.NonceSizeX)
	if n, _ := rand.Read(nonce); n != chacha20.NonceSizeX {
		return nil, errors.New("invalid nonce generation")
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(s.key, nonce)
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(nonce); err != nil {
		return nil, err
	}

	return &encryptWrite{file: f, cipher: cipher}, nil
}

type encryptRead struct {
	cipher *chacha20.Cipher
	file   io.ReadCloser
	buffer []byte
}

func (e *encryptRead) Read(p []byte) (n int, err error) {
	n, err = e.file.Read(p)
	if err != nil {
		return 0, err
	}

	if len(e.buffer) < n {
		e.buffer = make([]byte, n)
	}

	e.cipher.XORKeyStream(e.buffer, p[:n])
	return copy(p, e.buffer[:n]), err
}

func (e *encryptRead) Close() error {
	return e.file.Close()
}

type encryptWrite struct {
	cipher *chacha20.Cipher
	file   io.WriteCloser
	buffer []byte
}

func (e *encryptWrite) Write(p []byte) (n int, err error) {
	n = len(p)

	if len(e.buffer) < n {
		e.buffer = make([]byte, n)
	}

	e.cipher.XORKeyStream(e.buffer, p)
	return e.file.Write(e.buffer[:n])
}

func (e *encryptWrite) Close() error {
	return e.file.Close()
}
