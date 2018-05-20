package archivemanager

import "sort"

//BucketManager orchestrates the filling of buckets
type BucketManager struct {
	Buckets []Bucket
	Root    string
	Record  map[string]string
}

//InitializeBuckets generates buckets within the manager
func (m BucketManager) InitializeBuckets(numBuckets int, errChannel chan error) error {
	buckets, err := generateBuckets(m.Root, numBuckets, errChannel)
	if err != nil {
		return err
	}

	m.Buckets = buckets
	return nil
}

//CloseBuckets closes all writers and bucket channels
func (m BucketManager) CloseBuckets() error {
	for _, bucket := range m.Buckets {
		err := bucket.Writer.Close()
		if err != nil {
			return err
		}

		err = bucket.File.Close()
		if err != nil {
			return err
		}

		close(bucket.Channel)
	}
	return nil
}

func (m BucketManager) selectSmallestBucket() Bucket {
	sort.Slice(m.Buckets, func(i, j int) bool { return m.Buckets[i].Size < m.Buckets[j].Size })
	return m.Buckets[0]
}

//AddFileToBucket adds a file to the bucket
func (m BucketManager) AddFileToBucket(filename string) error {
	smallestBucket := m.selectSmallestBucket()
	err := smallestBucket.addToBucket(filename)
	if err != nil {
		return err
	}

	return nil
}
