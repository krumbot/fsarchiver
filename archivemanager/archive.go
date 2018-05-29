package archivemanager

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/krumbot/fsfileprocessor"
)

// Archive consumes the Controller options and starts the archiving process.
func Archive(crawlController fsfileprocessor.Controller, root string, numBuckets int) error {
	var compressedSize int64

	errChannel := make(chan error, 1)
	bm := BucketManager{
		Root: root,
	}

	err := bm.InitializeBuckets(numBuckets, errChannel)

	if err != nil {
		return err
	}

	processCb := func(walkinfo fsfileprocessor.WalkInfo) {
		bm.AddFileToBucket(walkinfo.Path)
	}

	go func() {
		err = <-errChannel
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	err = archive(crawlController, processCb)

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(bm.Record)
	if err != nil {
		return err
	}

	_, err = bm.RecordStore.Write(jsonData)
	if err != nil {
		return err
	}

	err = bm.CloseBuckets()

	if err != nil {
		return err
	}

	for _, bucket := range bm.Buckets {
		compressedSize += bucket.Size
	}

	fmt.Println("Compression Size: ", compressedSize)

	return nil
}

func archive(crawlController fsfileprocessor.Controller, processCb func(fsfileprocessor.WalkInfo)) error {

	crawlConfig := fsfileprocessor.Crawler{
		Processor:  generateProcessFunc(processCb),
		Controller: crawlController,
	}

	crawlErr := crawlConfig.Crawl()
	if crawlErr != nil {
		return crawlErr
	}

	return nil
}

func generateProcessFunc(processCb func(fsfileprocessor.WalkInfo)) func(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
	process := func(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
		for filewalkinfo := range fileReceiver {
			if !filewalkinfo.Info.IsDir() {
				processCb(filewalkinfo)
			}
		}
		close(errorChannel)
		return nil
	}

	return process
}

func cleanFile(filename string) error {
	return nil
}
