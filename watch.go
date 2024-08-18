package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var DRIVE_NAME string
var SOURCE_FOLDER string
var DESTINATION_FOLDER string
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
		if name == DRIVE_NAME {
			return id, name
		}
	}

	return "", ""
}

func Watch() {
	for range time.Tick(time.Second * time.Duration(INTERVAL)) {
		drive, name := detect()
		mounted := drive != ""
		if !MOUNTED && mounted {
			MOUNTED = mounted
			fmt.Printf("USB %s (%s:) injected\n", name, drive)
			filepath.WalkDir(fmt.Sprintf("%s:/%s", drive, SOURCE_FOLDER), visit)
			if COUNT != 0 {
				fmt.Printf("🍻 %d files copied 👏\n", COUNT)
			}
			COUNT = 0
		} else if MOUNTED && !mounted {
			MOUNTED = mounted
			fmt.Printf("USB %s (%s:) ejected\n", name, drive)
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

func visit(path string, entry fs.DirEntry, err error) error {
	filename := entry.Name()
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	extension := strings.ToUpper(filepath.Ext(filename))
	info, _ := entry.Info()
	if !entry.IsDir() && filter(extension) {
		yyyyMMdd := info.ModTime().Format("2006-01-02")
		dir := filepath.Join(DESTINATION_FOLDER, yyyyMMdd)
		instant := info.ModTime().Unix()
		p := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, instant, extension))

		// create dir per day
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create %s\n", dir)
			return nil
		}
		// move file to destination
		err = moveFile(path, p)
		if err != nil {
			fmt.Printf("Failed to copy %s to %s\n", path, p)
			return nil
		}
		fmt.Printf("Copied from %s to %s\n", path, p)
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
