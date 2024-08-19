cd cmd/cli
go build -o ../../usbcopy.exe

cd ../..

cd cmd/gui
go build -ldflags -H=windowsgui -o ../../usbcopy-gui.exe

cd ../..