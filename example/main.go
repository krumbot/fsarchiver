package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/krumbot/fsarchiver/archivemanager"
	"github.com/krumbot/fsfileprocessor"
)

func main() {
	err := exampleArchiver()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = exampleFileReader()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exampleFileReader() error {
	//Read from an existing bucket store
	bm := archivemanager.BucketManager{Root: "path/to/my/bucket-store"}

	err := bm.OpenExistingRecordStore()

	if err != nil {
		return err
	}

	//Retrieve an archived file from the bucket store
	buf, err := bm.RetrieveFile("path/of/the/file/to/retrieve")

	if err != nil {
		return err
	}

	fmt.Println(buf.String())

	return nil
}

func exampleArchiver() error {
	//Archive all json files that were modified bfore May 15, 2016 from the archive target
	fe, _ := regexp.Compile(".(json)")

	crawlController := fsfileprocessor.Controller{
		Rootdir:              "path/to/archive/target",
		Recursive:            true,
		EarliestTimeModified: time.Date(2016, time.May, 15, 0, 0, 0, 0, time.UTC),
		FileExt:              fe,
	}
	//And place them in evenly distributed buckets in the bucket-store directory
	err := archivemanager.Archive(crawlController, "path/to/my/bucket-store", 5)

	if err != nil {
		return err
	}

	return nil
}
