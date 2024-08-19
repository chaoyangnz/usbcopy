package main

import (
	"flag"
	"strings"
	"time"
	. "usbcopy/internal"
)

func main() {

	watchVolumePtr := flag.String("volume", "Z 6_2", "volume name to watch, no need to be full name")
	sourcePathPtr := flag.String("src", "%volume%:/DCIM", "absolute path to the source folder")
	destinationPathPtr := flag.String("dst", "E:/Photos/raw/%year%-%month%-%day%/%name%_%counter%.%extension%", "absolute path to the destination folder")
	filtersPtr := flag.String("filters", "NEF,JPEG,JPG,MOV,MP4", "a comma seperated list of file extensions without prefixed dot")
	intervalPtr := flag.Int("interval", 3, "watch interval in seconds")

	flag.Parse()

	context := &Context{
		UiMode:      false,
		WatchVolume: *watchVolumePtr,
		SrcPath:     *sourcePathPtr,
		DstPath:     *destinationPathPtr,
		Filters:     strings.Split(*filtersPtr, ","),
		Mounted:     false,
		Count:       0,
	}

	ticker := time.NewTicker(time.Second * time.Duration(*intervalPtr))
	Watch(ticker, context)
}
