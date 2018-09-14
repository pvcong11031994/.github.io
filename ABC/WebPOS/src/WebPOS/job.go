package main

import (
	"WebPOS/Common"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func startDeleteTmpFilesJob() {
	tmp_path, _ := filepath.Abs("./tmp")
	go func() {
		for true {
			old_dir := Common.DayFromToday(-2)

			folders, err := ioutil.ReadDir(tmp_path)

			if err == nil {
				for _, dir := range folders {
					if dir.IsDir() && dir.Name() <= old_dir {
						path := filepath.Join(tmp_path, dir.Name())
						err := os.RemoveAll(path)
						if err == nil {
							log.Println("RemoveAll " + path + " OK")
						} else {
							log.Println("RemoveAll " + path + " Fail")
						}
					}
				}
			}

			time.Sleep(time.Hour)
		}
	}()
}
