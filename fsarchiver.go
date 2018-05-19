package fsarchiver

import (
	"archive/zip"
	"io"
	"os"

	"github.com/krumbot/fsfileprocessor"
)

// Archive consumes the Controller options and starts the archiving process.
func Archive(crawlController fsfileprocessor.Controller) error {
	crawlConfig := fsfileprocessor.Crawler{
		Processor:  process,
		Controller: crawlController,
	}

	crawlErr := crawlConfig.Crawl()
	if crawlErr != nil {
		return crawlErr
	}

	return nil
}

// //Controller exposes a set of configuration options for the Archiver
// type Controller struct {
// 	ArchiveLocation string
// 	BundleArchives  bool
// }

func process(fileReceiver <-chan fsfileprocessor.WalkInfo, errorChannel chan<- error) error {
	for filewalkinfo := range fileReceiver {
		zipFile(filewalkinfo.Path, filewalkinfo.Path)
	}

	return nil
}

func zipFile(sourceFilename string, newFilename string) error {
	zipfile, err := os.Create(newFilename + ".zip")

	if err != nil {
		return err
	}
	defer zipfile.Close()

	writer := zip.NewWriter(zipfile)
	defer writer.Close()

	origFile, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer origFile.Close()

	info, err := origFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	writerHeader, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writerHeader, origFile)

	if err != nil {
		return err
	}

	return nil
}

func cleanFile(filename string) error {
	return nil
}
