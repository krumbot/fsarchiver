package archivemanager

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/satori/go.uuid"
)

//Bucket represents a zip sub-directory
type Bucket struct {
	path   string
	file   *os.File
	writer zip.Writer
	size   int64
}

//BucketManager orchestrates the filling of buckets
type BucketManager struct {
	buckets []Bucket
}

//GenerateStorageSystem is the main manager creator method
func GenerateStorageSystem(root string, numBuckets int) {
	buckets := generateBuckets(root, numBuckets)
}

func (m BucketManager) addFileToBucket(filename string) error {
	smallestBucket := m.selectSmallestBucket()
	err := addToBucket(smallestBucket, filename)
	if err != nil {
		return err
	}

	return nil
}

func (m BucketManager) selectSmallestBucket() Bucket {
	sort.Slice(m.buckets, func(i, j int) bool { return m.buckets[i].size < m.buckets[j].size })
	return m.Buckets[0]
}

func addToBucket(bucket Bucket, filename string) error {
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

	writer, err := bucket.writer.CreateHeader(fileHeader)
	if err != nil {
		return err
	}
	copySize, err = io.Copy(writer, fileToZip)
	if err != nil {
		return err
	}

	Bucket.size += copySize
	// write to json file
	return nil
}

func generateBuckets(rootPath string, num int) ([]Bucket, error) {
	buckets := make([]Bucket, num)
	for i := 0; i < num; i++ {
		newBucket, err := generateBucket(rootPath)
		if err != nil {
			return nil, err
		}

		buckets[i] = newBucket
	}

	return buckets, nil
}

func generateBucket(rootPath string) (Bucket, error) {
	hash, err := uuid.NewV4()
	if err != nil {
		return err
	}

	path := filepath.Join(rootPath, hash)
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := zip.NewWriter(file)

	newBucket := Bucket{path, file, writer}
	return newBucket, nil
}
