package gcs

import (
	"API-REST/services/conf"
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

type Storage struct {
	client  *storage.Client
	bucket  *storage.BucketHandle
	maxSize int
}

func Setup(s *Storage) error {
	var err error

	s = &Storage{
		maxSize: conf.Conf.GetInt("storage.gcs.maxSize"),
	}
	s.client, err = storage.NewClient(context.Background()) // Requires application_default_credentials.json configured in this system to work
	if err != nil {
		return err
	}
	s.bucket = s.client.Bucket(conf.Conf.GetString("storage.gcs.bucket"))

	return nil
}

func (s *Storage) UploadFile(filename string, b []byte) error {
	// Check file size
	size := len(b)
	if size > s.maxSize {
		return fmt.Errorf("file size too big: %d bytes (max=%d)", size, s.maxSize)
	}

	// Upload file
	w := s.bucket.Object(filename).NewWriter(context.TODO())
	_, err := w.Write(b)
	if err != nil {
		return err
	}
	return w.Close()
}

func (s *Storage) DownloadFile(filename string) ([]byte, error) {
	// Get file from gcs
	r, err := s.bucket.Object(filename).NewReader(context.TODO())
	if err != nil {
		return nil, err
	}
	// Get file content and return it
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Storage) DeleteFile(filename string) error {
	return s.bucket.Object(filename).Delete(context.TODO())
}
