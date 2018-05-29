# fsarchiver
### Still In development

This package uses channels to concurrently crawl through and archive a target directory.
Files are archived into evenly distributed buckets, which are then compressed. A lookup record is kept in the form of a .json file. Files can then be requested from the bucket store - the correct bucket will be unzipped and the file will be returned via a buffer.

See example/main.go for an example of archival and retrieval.

### TODO
- Add binary executable which takes a filepath as an argument and retrieves the file from the bucket store
- Add dummy archival artifact files in place of archival targets after archival process. These artifact files should link to the binary described above
- Add CLI options


