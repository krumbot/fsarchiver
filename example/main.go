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
	fe, _ := regexp.Compile(".(json)")

	crawlController := fsfileprocessor.Controller{
		Rootdir:              "your-src-dir",
		Recursive:            true,
		EarliestTimeModified: time.Date(2016, time.May, 15, 0, 0, 0, 0, time.UTC),
		FileExt:              fe,
	}

	err := archivemanager.Archive(crawlController, "your-output-dir", 5)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
