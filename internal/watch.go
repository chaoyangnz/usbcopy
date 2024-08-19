package internal

import (
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Context struct {
	// configs
	WatchVolume string
	UiMode      bool
	SrcPath     string
	DstPath     string
	Filters     []string
	// state
	Mounted bool
	Count   int
	// derived
	SrcBase    string
	VolumeId   string
	VolumeName string
}

func Watch(ticker *time.Ticker, contexts ...*Context) {
	for range ticker.C {
		for _, context := range contexts {
			id, name := detect(context.WatchVolume)
			context.VolumeId = id
			context.VolumeName = name
			mounted := name != ""
			if !context.Mounted && mounted {
				context.Mounted = mounted
				Run(context)
			} else if context.Mounted && !mounted {
				context.Mounted = mounted
				notify(context, "USB %s (%s:) ejected", name, id)
			} else {
				log("tick")
			}
		}
	}
}

func detect(namePattern string) (string, string) {
	args := []string{"logicaldisk", "where", "drivetype=2", "get", "deviceid,volumename"}
	cmd := exec.Command("wmic", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()

	if err != nil {
		return "", ""
	}

	s := string(out)

	l := strings.Split(s, "\r\n")

	if len(l) == 1 {
		return "", ""
	}

	for i := 0; i < len(l)-1; i++ {
		x := strings.Split(strings.TrimSpace(l[i+1]), ":")
		if len(x) < 2 {
			return "", ""
		}
		id := x[0]
		name := strings.TrimSpace(x[1])
		if strings.Contains(name, namePattern) {
			return id, name
		}
	}

	return "", ""
}
