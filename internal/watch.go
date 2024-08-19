package internal

import (
	"time"
)

type Context struct {
	// configs
	UiMode  bool
	Volume  string
	SrcPath string
	DstPath string
	Filters []string
	// state
	Mounted bool
	Count   int
	// derived
	SrcBase string
}

func Watch(ticker *time.Ticker, contexts ...*Context) {
	for range ticker.C {
		for _, context := range contexts {
			run(context)
		}
	}
}
