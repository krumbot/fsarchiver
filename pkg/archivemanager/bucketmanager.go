package archivemanager

import (
	"os"
	"path/filepath"
	"sort"
)

//BucketManager orchestrates the filling of buckets
type BucketManager struct {
	Buckets     []Bucket
	Root        string
	Record      map[string]string
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

	m.Record = make(map[string]string)

	return nil
}

func (m *BucketManager) initializeRecordStore() error {
	file, err := os.Create(filepath.Join(m.Root, ".bucket-store.json"))

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
