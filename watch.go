package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var VOLUME_NAME string
var SOURCE_PATH string
var DESTINATION_PATH string
var FILE_EXTENSIONS []string = nil
var INTERVAL = 3

var MOUNTED = false
var COUNT = 0

func detect() (string, string) {
	//var drives []string
	//driveMap := make(map[string]bool)

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
		if strings.Contains(name, VOLUME_NAME) {
			return id, name
		}
	}

	return "", ""
}

func Watch() {
	for range time.Tick(time.Second * time.Duration(INTERVAL)) {
		volume, name := detect()
		sourcebase := interpolate(SOURCE_PATH, []string{
			"%volume%",
		}, []string{
			volume,
		})
		mounted := volume != ""
		if !MOUNTED && mounted {
			MOUNTED = mounted
			fmt.Printf("USB %s (%s:) injected\n", name, volume)
			filepath.WalkDir(sourcebase, visit)
			if COUNT != 0 {
				fmt.Printf("🍻 %d files copied 👏\n", COUNT)
			}
			COUNT = 0
		} else if MOUNTED && !mounted {
			MOUNTED = mounted
			fmt.Printf("USB %s (%s:) ejected\n", name, volume)
		} else {
			fmt.Printf("tick\n")
		}
	}
}

func filter(extension string) bool {
	for _, ext := range FILE_EXTENSIONS {
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

func getExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.ToUpper(ext[1:])
}

func visit(sourcepath string, entry fs.DirEntry, err error) error {

	sourcedir := filepath.Dir(sourcepath)
	volume := filepath.VolumeName(sourcepath)
	filename := entry.Name()
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	extension := getExtension(filename)
	info, _ := entry.Info()

	if !entry.IsDir() && filter(extension) {
		year := info.ModTime().Format("2006")
		month := info.ModTime().Format("01")
		day := info.ModTime().Format("02")
		midnight, _ := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", year, month, day))
		diff := strconv.FormatInt(info.ModTime().Unix()-midnight.Unix(), 10)
		sourcebase := interpolate(SOURCE_PATH, []string{
			"%volume%",
		}, []string{
			volume,
		})
		sourcereldir, _ := filepath.Rel(sourcebase, sourcedir)
		destpath := interpolate(DESTINATION_PATH, []string{
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
			sourcereldir,
			filename,
			name,
			extension,
			diff,
		})
		destdir := filepath.Dir(destpath)
		// create dir first
		err := os.MkdirAll(destdir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create %s\n", destdir)
			return nil
		}
		// move file to destination
		err = moveFile(sourcepath, destpath)
		if err != nil {
			fmt.Printf("Failed to copy %s to %s %v\n", sourcepath, destpath, err)
			return nil
		}
		fmt.Printf("Copied from %s to %s\n", sourcepath, destpath)
		COUNT += 1
	}

	return nil
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
