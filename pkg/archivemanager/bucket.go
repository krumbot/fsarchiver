package archivemanager

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/satori/go.uuid"
)

//Bucket represents a zip sub-directory
type Bucket struct {
	Path   string
	File   *os.File
	Writer *zip.Writer
	Size   int64
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

	infoClone := fileInfo
	err = copier.Copy(&infoClone, &fileInfo)

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

	bm.Record[hash.String()] = filename

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
	file, err := os.Create(path + ".zip")
	if err != nil {
		return Bucket{}, err
	}

	writer := zip.NewWriter(file)

	newBucket := Bucket{path, file, writer, 0}

	return newBucket, nil
}
