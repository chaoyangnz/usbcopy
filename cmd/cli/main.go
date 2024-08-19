package main

import (
	"flag"
	"log"
	"os"
	"strings"
	. "usbcopy/internal"
)

func main() {

	volumeIdPtr := flag.String("volume", "Z", "volume letter")
	sourcePathPtr := flag.String("src", "%volume%:/DCIM", "absolute path in source drive")
	destinationPathPtr := flag.String("dst", "E:/Photos/raw/%year%-%month%-%day%/%name%_%counter%.%extension%", "absolute path to the destination folder")
	filtersPtr := flag.String("filters", "NEF,JPEG,JPG,MOV,MP4", "a comma seperated list of file extensions without prefixed dot")

	flag.Parse()

	context := &Context{
		UiMode:   true,
		VolumeId: *volumeIdPtr,
		SrcPath:  *sourcePathPtr,
		DstPath:  *destinationPathPtr,
		Filters:  strings.Split(*filtersPtr, ","),
		Mounted:  false,
		Count:    0,
	}

	// set logging
	f, err := os.OpenFile("usbcopy.log", os.O_RDWR|os.O_CREATE, 0666)
	if err == nil {
		log.SetOutput(f)
	}

	Run(context)

}
