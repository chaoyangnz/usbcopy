package internal

import (
	"fmt"
	"github.com/go-toast/toast"
	"io"
	"io/fs"
	logging "log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

func detect(context *Context) (string, string) {
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
		if strings.Contains(name, context.Volume) {
			return id, name
		}
	}

	return "", ""
}

func Watch(ticker *time.Ticker, contexts ...*Context) {
	for range ticker.C {
		for _, context := range contexts {
			volume, name := detect(context)
			sourcebase := interpolate(context.SrcPath, []string{
				"%volume%",
			}, []string{
				volume,
			})
			context.SrcBase = sourcebase
			mounted := volume != ""
			if !context.Mounted && mounted {
				context.Mounted = mounted
				notify(context, "USB %s (%s:) injected", name, volume)
				filepath.WalkDir(sourcebase, visitFn(context))
				if context.Count != 0 {
					notify(context, "🍻 %d files copied 👏", context.Count)
				}
				context.Count = 0
			} else if context.Mounted && !mounted {
				context.Mounted = mounted
				notify(context, "USB %s (%s:) ejected", name, context.Volume)
			} else {
				log("tick")
			}
		}
	}
}

func visitFn(context *Context) fs.WalkDirFunc {
	return func(sourcepath string, entry fs.DirEntry, err error) error {
		sourcedir := filepath.Dir(sourcepath)
		filename := entry.Name()
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		extension := extractExtension(filename)
		info, _ := entry.Info()

		if !entry.IsDir() && filter(extension, context.Filters) {
			year := info.ModTime().Format("2006")
			month := info.ModTime().Format("01")
			day := info.ModTime().Format("02")
			midnight, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))
			diff := strconv.FormatInt(info.ModTime().Unix()-midnight.Unix(), 10)
			sourcedirRel, _ := filepath.Rel(context.SrcBase, sourcedir)
			destpath := interpolate(context.DstPath, []string{
				"%year%",
				"%month%",
				"%day%",
				"%dir%",
				"%filename%",
				"%name%",
				"%extension%",
				"%counter%",
			}, []string{
				year,
				month,
				day,
				sourcedirRel,
				filename,
				name,
				extension,
				diff,
			})
			destdir := filepath.Dir(destpath)
			// create dir first
			err := os.MkdirAll(destdir, os.ModePerm)
			if err != nil {
				log("Failed to create %s", destdir)
				return nil
			}
			// move file to destination
			err = moveFile(sourcepath, destpath)
			if err != nil {
				log("Failed to copy %s to %s %v", sourcepath, destpath, err)
				return nil
			}
			log("Copied from %s to %s", sourcepath, destpath)
			context.Count += 1
		}

		return nil
	}
}

func log(format string, args ...any) {
	logging.Printf(format+"\n", args...)
}

func notify(context *Context, format string, args ...any) {
	if context.UiMode {
		notification := toast.Notification{
			AppID:   "usbcopy",
			Title:   "🛈",
			Message: fmt.Sprintf(format, args...),
		}
		err := notification.Push()
		if err != nil {
			log("%v", err)
		}
	} else {
		log(format, args...)
	}
}

func filter(extension string, whitelist []string) bool {
	for _, ext := range whitelist {
		if ext == extension {
			return true
		}
	}
	return false
}

func interpolate(s string, vars []string, vals []string) string {
	str := s
	for i, _ := range vars {
		str = strings.ReplaceAll(str, vars[i], vals[i])
	}
	return filepath.Clean(str)
}

func extractExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.ToUpper(ext[1:])
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}
