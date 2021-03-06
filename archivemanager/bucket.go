package archivemanager

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

//Bucket represents a zip sub-directory
type Bucket struct {
	Name   string
	File   *os.File
	Writer *zip.Writer
	Size   int64
	Reader *zip.ReadCloser
}

func (bucket *Bucket) addToBucket(filename string, bm *BucketManager) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer fileToZip.Close()

	fileInfo, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	hash, err := uuid.NewV4()
	if err != nil {
		return err
	}

	fileHeader := zip.FileHeader{
		Name:     hash.String(),
		Method:   zip.Deflate,
		Modified: fileInfo.ModTime(),
	}

	writer, err := bucket.Writer.CreateHeader(&fileHeader)
	if err != nil {
		return err
	}
	copySize, err := io.Copy(writer, fileToZip)
	if err != nil {
		return err
	}

	bm.Record[bucket.Name][filename] = hash.String()

	bucket.Size += copySize

	return nil
}

func openExistingBucket(root string, name string) (Bucket, error) {
	var size int64

	path := filepath.Join(root, name) + ".zip"
	file, err := os.Open(path)

	if err != nil {
		return Bucket{}, err
	}

	reader, err := zip.OpenReader(path)

	if err != nil {
		return Bucket{}, err
	}

	for _, f := range reader.File {
		size += int64(f.UncompressedSize64)
	}

	bucket := Bucket{Name: name, File: file, Size: size, Reader: reader}
	return bucket, nil
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
	file, err := os.Create(path + ".zip")
	if err != nil {
		return Bucket{}, err
	}

	writer := zip.NewWriter(file)

	reader, err := zip.OpenReader(path)

	if err != nil {
		return Bucket{}, err
	}

	newBucket := Bucket{hash.String(), file, writer, 0, reader}

	return newBucket, nil
}

func (bucket *Bucket) fetchFile(lookup string) (io.ReadCloser, error) {
	for _, f := range bucket.Reader.File {
		if f.Name == lookup {
			fileContent, err := f.Open()

			if err != nil {
				return nil, err
			}
			return fileContent, nil
		}
	}
	return nil, errors.New("Could not find file with index " + lookup + "in bucket " + bucket.Name)
}
