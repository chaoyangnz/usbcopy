cd cmd/cli
go build -ldflags="-H windowsgui" -o ../../usbcopy.exe

cd ../..

cd cmd/cli-watch
go build -o ../../usbcopy.watch.exe

cd ../..

cd cmd/gui
go build -ldflags="-H windowsgui" -o ../../usbcopy.gui.exe

cd ../..

del E:/usbcopy*
copy *.exe E:
copy *.yml E: