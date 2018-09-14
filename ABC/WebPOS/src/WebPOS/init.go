package main

import (
	"WebPOS/Common"
	"github.com/goframework/gf/ext"
	"github.com/goframework/gf/exterror"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const _WINDOWS_OS = "windows"
const _LOG_FILE_SPLITTER = "RELEASE"
const _PATH_LOG = "../log/WEBPOS/"

func init() {

	// Begin log console debug
	//log.SetOutput(os.Stdout)
	//log.Println("Setting output log into stdout")

	logFileExec, _ := filepath.Abs(os.Args[0])
	log.Println("Path file excecute is " + logFileExec)

	log.Println("Create file log - START")
	if runtime.GOOS != _WINDOWS_OS {
		//logFilePath := _PATH_LOG + "WebPOSLinux.exe.log"
		logFilePath := _PATH_LOG + "application.log"
		log.Println("Setting log filename: " + logFilePath)

		log.Println("Check path 「" + _PATH_LOG + "」 to save file log - START")
		if !ext.FolderExists(_PATH_LOG) {
			log.Println("In case of path not exists")
			os.MkdirAll(_PATH_LOG, os.ModePerm)

			log.Println("Create path: " + _PATH_LOG)
			// Check if folder is created without error
			if ext.FolderExists(_PATH_LOG) {
				log.Println("Created at " + _PATH_LOG + " success")
			} else {
				log.Println("Created at " + _PATH_LOG + " fail")
			}
		} else {
			log.Println("Path exist at " + _PATH_LOG)
		}
		log.Println("Check path 「" + _PATH_LOG + "」to save file log - END")

		log.Println("Start create file log - START")
		f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		log.Println("Start create file log - Success")
		log.Println("Start create file log - END")

		Common.LogErr(err)
		log.Println("Write log sample to file log")

		// Check if file is created
		log.Println("Check exists file log: " + logFilePath + " - START")
		_, err = os.Stat(logFilePath)
		if os.IsNotExist(err) {
			log.Println("Log file not exist at " + logFilePath)
		} else {
			log.Println("Log file exist at " + logFilePath)
		}
		log.Println("Check exists file log: " + logFilePath + " - END")

		// End log console debug, logged back to file
		log.Println("Create file log - END")
		log.SetOutput(f)
	}

	exterror.SetFilePathSpliter(_LOG_FILE_SPLITTER)

	startDeleteTmpFilesJob()
}
