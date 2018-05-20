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
	fe, _ := regexp.Compile(".(md)")

	crawlController := fsfileprocessor.Controller{
		Rootdir:              "../fsfileprocessor/",
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

	processCb := func(walkinfo fsfileprocessor.WalkInfo) error {
		err := bm.AddFileToBucket(walkinfo.Path)
		if err != nil {
			return err
		}
		return nil
	}
}
