package internal

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func run(context *Context) {
	volume, name := detect(context)
	sourcebase := interpolate(context.SrcPath, []string{
		"%volume%",
	}, []string{
		volume,
	})
	mounted := volume != ""
	if !context.Mounted && mounted {
		context.Mounted = mounted
		context.SrcBase = sourcebase
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

func visitFn(context *Context) fs.WalkDirFunc {
	return func(srcPath string, entry fs.DirEntry, err error) error {
		srcDir := filepath.Dir(srcPath)
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
			srcDirRel, _ := filepath.Rel(context.SrcBase, srcDir)
			dstPath := interpolate(context.DstPath, []string{
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
				srcDirRel,
				filename,
				name,
				extension,
				diff,
			})
			destDir := filepath.Dir(dstPath)
			// create dir first
			err := os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				log("Failed to create %s", destDir)
				return nil
			}
			// move file to destination
			err = moveFile(srcPath, dstPath)
			if err != nil {
				log("Failed to copy %s to %s %v", srcPath, dstPath, err)
				return nil
			}
			log("Copied from %s to %s", srcPath, dstPath)
			context.Count += 1
		}

		return nil
	}
}
