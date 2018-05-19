package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/krumbot/fsarchiver"
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

	archiveErr := fsarchiver.Archive(crawlController)

	if archiveErr != nil {
		fmt.Println(archiveErr)
		os.Exit(1)
	}
}
