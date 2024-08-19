package main

import (
	"flag"
	"strings"
	"time"
	. "usbcopy/internal"
)

func main() {

	volumeNamePtr := flag.String("volume", "Z 6_2", "volume name")
	sourcePathPtr := flag.String("src-path", "%volume%:/DCIM", "absolute path in source drive")
	destinationPathPtr := flag.String("dst-path", "E:/Photos/raw/%year%-%month%-%day%/%name%_%counter%.%extension%", "absolute path to the destination folder")
	filtersPtr := flag.String("filters", "NEF,JPEG,JPG,MOV,MP4", "a comma seperated list of file extensions without prefixed dot")
	intervalPtr := flag.Int("interval", 3, "watch interval in seconds")

	flag.Parse()

	context := &Context{
		UiMode:  false,
		Volume:  *volumeNamePtr,
		SrcPath: *sourcePathPtr,
		DstPath: *destinationPathPtr,
		Filters: strings.Split(*filtersPtr, ","),
		Mounted: false,
		Count:   0,
	}

	ticker := time.NewTicker(time.Second * time.Duration(*intervalPtr))

	Watch(ticker, context)
}
