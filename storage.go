package giostorage

import (
	"errors"
	"golang.org/x/crypto/chacha20"
	"io"
)

var (
	EncryptedKeySize = chacha20.KeySize
)

// NewPlainStorage creates a Storage.
func NewPlainStorage(name string) *Storage {
	return &Storage{storage: &plainStorage{name: name}}
}

// NewEncryptedStorage creates a Storage encrypted by unauthenticated ChaCha20X.
// The key must be chacha20.KeySize. The salt is created at random, using rand.Read.
//
// The file can be modified without any sign of change, similarly as plain-text files.
// Beware that the key is already on the client's computer, so it's possible to
// re-encrypt an arbitrary data.
func NewEncryptedStorage(name string, key []byte) (*Storage, error) {
	if len(key) != EncryptedKeySize {
		return nil, errors.New("invalid key size")
	}

	return &Storage{storage: &encryptedStorage{plain: &plainStorage{name: name}, key: key}}, nil
}

// Storage provides the Open and Create functions, using an arbitrary storage implementation.
// You should create the Storage from NewPlainStorage or NewEncryptedStorage.
type Storage struct {
	// storage is private to avoid leak methods, like .PlainStorage.Open while using EncryptedStorage.
	storage storage
}

type storage interface {
	open(name string) (io.ReadCloser, error)
	create(name string) (io.WriteCloser, error)
}

// Open opens the file.
func (s *Storage) Open(name string) (io.ReadCloser, error) {
	return s.storage.open(name)
}

// Create creates (or truncate, if exists) the file.
func (s *Storage) Create(name string) (io.WriteCloser, error) {
	return s.storage.create(name)
}

// ReadInterface will read the given resp interface, from the file.
// It uses gob or json.
func (s *Storage) ReadInterface(name string, resp interface{}) error {
	f, err := s.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Decoder(f).Decode(resp)
}

// WriteInterface will write given data interface into a file.
// It uses gob or json.
func (s *Storage) WriteInterface(name string, data interface{}) error {
	f, err := s.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Encoder(f).Encode(data)
}
