package main

import (
	_ "embed"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
	"os"
	"strings"
	"time"
	. "usbcopy/internal"
)

type VolumeConfig struct {
	Volume  string `mapstructure:"volume"`
	SrcPath string `mapstructure:"src"`
	DstPath string `mapstructure:"dst"`
	Filters string `mapstructure:"filters"`
}

type Config struct {
	Interval int            `mapstructure:"interval"`
	Volumes  []VolumeConfig `mapstructure:"volumes"`
}

func main() {
	var conf Config = Config{}

	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("usbcopy.yml")

	if err == nil {
		config.Decode(&conf)
	} else {
		conf = Config{
			Volumes: []VolumeConfig{
				{
					Volume:  "Z 6_2",
					SrcPath: "%volume%:/DCIM",
					DstPath: "E:/Photos/raw/%year%-%month%-%day%/%name%_%counter%.%extension%",
					Filters: "NEF,JPEG,JPG,MOV,MP4",
				},
			},
		}
	}

	contexts := make([]*Context, len(conf.Volumes))

	for i, conf := range conf.Volumes {
		contexts[i] = &Context{
			UiMode:      true,
			WatchVolume: conf.Volume,
			SrcPath:     conf.SrcPath,
			DstPath:     conf.DstPath,
			Filters:     strings.Split(conf.Filters, ","),
			Mounted:     false,
			Count:       0,
		}
	}

	// set logging
	f, err := os.OpenFile("usbcopy.log", os.O_RDWR|os.O_CREATE, 0666)
	if err == nil {
		log.SetOutput(f)
	}

	ticker := time.NewTicker(time.Second * time.Duration(config.Int("interval", 3)))

	systray.Run(onReadyFn(ticker, contexts), onExitFn(ticker, contexts))
}

//go:embed icon.ico
var icon []byte

func onReadyFn(ticker *time.Ticker, contexts []*Context) func() {
	return func() {
		systray.SetIcon(icon)
		systray.SetTitle("usbcopy")
		systray.SetTooltip("copy usb to destination files automatically")
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		// Sets the icon of a menu item. Only available on Mac and Windows.
		mQuit.SetIcon(icon)

		go func() {
			<-mQuit.ClickedCh
			fmt.Println("Requesting quit")
			systray.Quit()
			fmt.Println("Finished quitting")
		}()

		go func() {
			Watch(ticker, contexts...)
		}()
	}
}

func onExitFn(ticker *time.Ticker, contexts []*Context) func() {
	return func() {
		// clean up here
		ticker.Stop()
	}
}
