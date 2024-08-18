# usbcopy

A utility to watch USB/Removable drives and copy files to the destination folder automatically. This will largely improve you photography post-processing workflow.

Inspired by Synology USB Copy add-on package.

## How to use

```
C:\> usbcopy.exe -h

Usage of usbcopy.exe:
  -destination-folder string
        absolute path to the destination folder (default "E:/Photos/raw")
  -drive-name string
        drive name (default "NIKON Z 6_2")
  -file-extensions string
        a comma seperated list of file extensions (default ".NEF,.JPEG,.JPG,.MOV,.MP4")
  -interval int
        watch interval (default 3)
  -source-folder string
        relative path in source drive (default "DCIM")

```

## Roadmap

- [ ] Windows service wrapper
- [ ] Windows desktop notification integration
- [ ] better logging & configuration
- [ ] flexible file renaming strategies and directory layout 