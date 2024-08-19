package main

import (
	"flag"
	"strings"
)

func main() {
	volumeNamePtr := flag.String("volume-name", "Z 6_2", "volume name")
	sourcePathPtr := flag.String("source-path", "%volume%:/DCIM", "absolute path in source drive")
	destinationPathPtr := flag.String("destination-path", "E:/Photos/raw/%year%-%month%-%day%/%name%_%counter%.%extension%", "absolute path to the destination folder")
	fileExtensionsPtr := flag.String("file-extensions", "NEF,JPEG,JPG,MOV,MP4", "a comma seperated list of file extensions without prefixed dot")
	intervalPtr := flag.Int("interval", 3, "watch interval in seconds")

	flag.Parse()

	VOLUME_NAME = *volumeNamePtr
	SOURCE_PATH = *sourcePathPtr
	DESTINATION_PATH = *destinationPathPtr
	FILE_EXTENSIONS = strings.Split(*fileExtensionsPtr, ",")
	INTERVAL = *intervalPtr

	Watch()
}
