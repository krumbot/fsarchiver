package archivemanager

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

//Bucket represents a zip sub-directory
type Bucket struct {
	Path    string
	File    *os.File
	Writer  *zip.Writer
	Size    int64
	Channel chan string
}

func (bucket Bucket) addToBucket(filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer fileToZip.Close()

	fileInfo, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	fileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	fileHeader.Method = zip.Deflate

	writer, err := bucket.Writer.CreateHeader(fileHeader)
	if err != nil {
		return err
	}
	copySize, err := io.Copy(writer, fileToZip)
	if err != nil {
		return err
	}

	bucket.Size += copySize
	return nil
}

func generateBuckets(rootPath string, num int, errChannel chan error) ([]Bucket, error) {
	buckets := make([]Bucket, num)
	for i := 0; i < num; i++ {
		newBucket, err := generateBucket(rootPath, errChannel)
		if err != nil {
			return nil, err
		}

		buckets[i] = newBucket
	}

	return buckets, nil
}

func generateBucket(rootPath string, errChannel chan error) (Bucket, error) {
	hash, err := uuid.NewV4()
	if err != nil {
		return Bucket{}, err
	}

	path := filepath.Join(rootPath, hash.String())
	file, err := os.Create(path)
	if err != nil {
		return Bucket{}, err
	}

	writer := zip.NewWriter(file)

	channel := make(chan string)
	newBucket := Bucket{path, file, writer, 0, channel}

	go func() {
		for filename := range channel {
			err := newBucket.addToBucket(filename)
			if err != nil {
				errChannel <- err
			}
		}
		close(channel)
	}()

	return newBucket, nil
}
