package local

import (
	"API-REST/services/conf"
	"fmt"
	"io/ioutil"
	"os"
)

type Storage struct {
	storageDir string
	maxSize    int
}

func Setup(s *Storage) error {
	s = &Storage{
		storageDir: conf.Conf.GetString("storage.local.rootDir"),
		maxSize:    conf.Conf.GetInt("storage.local.maxSize"),
	}

	return nil
}

func (s *Storage) SaveFile(filename string, b []byte) error {
	// Check file size
	size := len(b)
	if size > s.maxSize {
		return fmt.Errorf("file size too big: %d bytes (max=%d)", size, s.maxSize)
	}

	// open output file
	out, err := os.Create(s.storageDir + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// write content into output file
	_, err = out.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetFile(filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(s.storageDir + "/" + filename)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Storage) DeleteFile(filename string) error {
	return os.Remove(s.storageDir + "/" + filename)
}
