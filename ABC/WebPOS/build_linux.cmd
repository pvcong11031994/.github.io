@echo off

@SET GOOS=linux
@SET GOARCH=amd64
@SET GOROOT=C:/Go
@SET GOPATH=%CD%;%CD%/../GOLIB
@SET FLAG=-ldflags "-s -w"
@SET OUT_DIR=BUILD\Linux
@SET OUTPUT_FILE_NAME=WebPOSLinux
@SET EXT=.exe
@SET CMD=go build

@SET BUILD_PROGRAM=WebPOS

echo Building %BUILD_PROGRAM% ...
%CMD% %FLAG% -o %OUT_DIR%\%OUTPUT_FILE_NAME%%EXT% %BUILD_PROGRAM%