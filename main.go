package main

import (
	"flag"
	"strings"
)

func main() {
	driveNamePtr := flag.String("drive-name", "NIKON Z 6_2", "drive name")
	sourceFolderPtr := flag.String("source-folder", "DCIM", "relative path in source drive")
	destinationFolderPtr := flag.String("destination-folder", "E:/Photos/raw", "absolute path to the destination folder")
	fileExtensionsPtr := flag.String("file-extensions", ".NEF,.JPEG,.JPG,.MOV,.MP4", "a comma seperated list of file extensions")
	intervalPtr := flag.Int("interval", 3, "watch interval")

	flag.Parse()

	DRIVE_NAME = *driveNamePtr
	SOURCE_FOLDER = *sourceFolderPtr
	DESTINATION_FOLDER = *destinationFolderPtr
	FILE_EXTENSIONS = strings.Split(*fileExtensionsPtr, ",")
	INTERVAL = *intervalPtr

	Watch()
}
