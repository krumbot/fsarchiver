package archivemanager

import (
	"github.com/krumbot/fsfileprocessor"
)

// Archive consumes the Controller options and starts the archiving process.
func Archive(crawlController fsfileprocessor.Controller, processCb func(fsfileprocessor.WalkInfo) error) error {

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

func generateProcessFunc(processCb func(fsfileprocessor.WalkInfo) error) func(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
	process := func(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
		for filewalkinfo := range fileReceiver {
			err := processCb(filewalkinfo)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return process
}

func cleanFile(filename string) error {
	return nil
}
