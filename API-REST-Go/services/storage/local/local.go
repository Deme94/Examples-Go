package local

import (
	"API-REST/services/conf"
	"fmt"
	"io"
	"os"
)

type Storage struct {
	storageDir string
	maxSize    int64
}

func Setup(s *Storage) error {
	s = &Storage{
		storageDir: conf.Conf.GetString("storage.local.rootDir"),
		maxSize:    conf.Conf.GetInt64("storage.local.maxSize"),
	}

	return nil
}

func (s *Storage) SaveFile(filename string, f *os.File) error {
	// Resets read offset, just in case
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check file size
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	if size > s.maxSize {
		return fmt.Errorf("file size too big: %d bytes (max=%d)", size, s.maxSize)
	}

	// open output file
	out, err := os.Create(s.storageDir + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		_, err = out.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetFile(filename string) (*os.File, error) {
	f, err := os.Open(s.storageDir + "/" + filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Storage) DeleteFile(filename string) error {
	return os.Remove(s.storageDir + "/" + filename)
}
