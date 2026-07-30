package media

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("media: not found")

type Storage interface {
	Save(context.Context, string, io.Reader) error
	Delete(context.Context, string) error
}

type LocalStorage struct{ root string }

func NewLocalStorage(root string) *LocalStorage {
	return &LocalStorage{root: root}
}

func (s *LocalStorage) Save(_ context.Context, key string, source io.Reader) error {
	if err := os.MkdirAll(s.root, 0o750); err != nil {
		return err
	}
	target, err := os.OpenFile(filepath.Join(s.root, key), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		return err
	}
	defer target.Close()
	_, err = io.Copy(target, source)
	return err
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	err := os.Remove(filepath.Join(s.root, key))
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	return err
}
