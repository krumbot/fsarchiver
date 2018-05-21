package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/krumbot/fsarchiver/pkg/archivemanager"
	"github.com/krumbot/fsfileprocessor"
)

func main() {
	fe, _ := regexp.Compile(".(xml)")

	crawlController := fsfileprocessor.Controller{
		Rootdir:              "/Users/vikrum/Documents/zipsrc",
		Recursive:            true,
		EarliestTimeModified: time.Date(2016, time.May, 15, 0, 0, 0, 0, time.UTC),
		FileExt:              fe,
	}

	errChannel := make(chan error, 1)
	bm := archivemanager.BucketManager{
		Root: "/Users/vikrum/Documents/ziptest/",
	}

	err := bm.InitializeBuckets(2, errChannel)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	err = archivemanager.Archive(crawlController, processCb)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = bm.CloseBuckets()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, bucket := range bm.Buckets {
		fmt.Println(int(bucket.Size) / 1000)
	}

}
