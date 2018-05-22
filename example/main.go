package main

import (
	"fmt"

	"github.com/krumbot/fsarchiver/pkg/archivemanager"
)

func main() {

	bm := archivemanager.BucketManager{Root: "your-output-path"}

	bm.OpenExistingRecordStore()

	for _, bucket := range bm.Buckets {
		fmt.Println(bucket)
	}

	// fe, _ := regexp.Compile(".(json)")

	// crawlController := fsfileprocessor.Controller{
	// 	Rootdir:              "your-src-path",
	// 	Recursive:            true,
	// 	EarliestTimeModified: time.Date(2016, time.May, 15, 0, 0, 0, 0, time.UTC),
	// 	FileExt:              fe,
	// }

	// err := archivemanager.Archive(crawlController, "your-output-path", 5)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

}
