package archivemanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

const bucketStoreFilename = ".bucket-store.json"

//BucketManager orchestrates the filling of buckets
type BucketManager struct {
	Buckets     []Bucket
	Root        string
	Record      map[string]map[string]string
	RecordStore *os.File
}

//InitializeBuckets generates buckets within the manager
func (m *BucketManager) InitializeBuckets(numBuckets int, errChannel chan error) error {
	buckets, err := generateBuckets(m.Root, numBuckets, errChannel)
	if err != nil {
		return err
	}

	m.Buckets = buckets

	err = m.initializeRecordStore()
	if err != nil {
		return nil
	}

	m.Record = make(map[string]map[string]string)

	for _, bucket := range m.Buckets {
		m.Record[bucket.Name] = make(map[string]string)
	}

	return nil
}

//RetrieveFile is meant to find the correct bucket for a file, fetch it, and return it
func (m *BucketManager) RetrieveFile(filename string) (*bytes.Buffer, error) {
	bucket, lookup, err := m.fileBucketLookup(filename)
	if err != nil {
		return nil, err
	}

	fileContent, err := bucket.fetchFile(lookup)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(fileContent)
	return buf, nil
}

func (m *BucketManager) fileBucketLookup(filename string) (Bucket, string, error) {
	for bucketName := range m.Record {
		lookupValue := m.Record[bucketName][filename]
		if lookupValue != "" {
			bucket, err := m.fetchBucket(bucketName)
			if err != nil {
				return Bucket{}, "", err
			}

			return bucket, lookupValue, nil
		}
	}

	return Bucket{}, "", errors.New("No bucket was found containing file: " + filename)
}

func (m *BucketManager) fetchBucket(bucketName string) (Bucket, error) {
	for _, bucket := range m.Buckets {
		if bucket.Name == bucketName {
			return bucket, nil
		}
	}
	return Bucket{}, errors.New("Invalid bucket lookup by bucket name: bucket " + bucketName + "was not found.")
}

//OpenExistingRecordStore opens existing bucket store data
func (m *BucketManager) OpenExistingRecordStore() error {
	bucketStore := filepath.Join(m.Root, bucketStoreFilename)

	record := make(map[string]map[string]string)

	jsonContent, err := ioutil.ReadFile(bucketStore)

	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonContent, &record)
	if err != nil {
		return err
	}

	recordStore, err := os.Open(bucketStore)
	if err != nil {
		return err
	}

	m.Record = record
	m.RecordStore = recordStore

	for k := range record {
		bucket, err := openExistingBucket(m.Root, k)
		if err != nil {
			return err
		}

		m.Buckets = append(m.Buckets, bucket)
	}

	return nil

}

func (m *BucketManager) initializeRecordStore() error {
	file, err := os.Create(filepath.Join(m.Root, bucketStoreFilename))

	if err != nil {
		return err
	}

	m.RecordStore = file
	return nil
}

//CloseBuckets closes all writers
func (m *BucketManager) CloseBuckets() error {
	for _, bucket := range m.Buckets {
		err := bucket.Writer.Close()
		if err != nil {
			return err
		}

		err = bucket.File.Close()
		if err != nil {
			return err
		}
	}

	err := m.RecordStore.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *BucketManager) selectSmallestBucket() *Bucket {
	sort.Slice(m.Buckets, func(i, j int) bool { return m.Buckets[i].Size < m.Buckets[j].Size })
	return &m.Buckets[0]
}

//AddFileToBucket adds a file to the bucket
func (m *BucketManager) AddFileToBucket(filename string) {
	smallestBucket := m.selectSmallestBucket()
	smallestBucket.addToBucket(filename, m)
}
