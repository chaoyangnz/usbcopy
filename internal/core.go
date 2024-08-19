package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Run(context *Context) {
	notify(context, "USB %s: found %s", context.VolumeId, context.VolumeName)
	srcBase := interpolate(context.SrcPath, []string{
		"%volume%",
	}, []string{
		context.VolumeId,
	})
	context.SrcBase = srcBase
	filepath.WalkDir(srcBase, visitFn(context))
	if context.Count != 0 {
		notify(context, "🍻 %d files copied 👏", context.Count)
	}
	context.Count = 0
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
