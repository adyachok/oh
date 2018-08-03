package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"os"

	"github.com/adyachok/oh/remote"
	"github.com/adyachok/oh/web"
	"github.com/adyachok/oh/models"
)

// Interface defines method to get information about files
type dirInfo interface {
	listDir() []os.FileInfo
}

// Local implementation of the dirInfo interface.
// Parses current dirrectory and builds list of FileInfo objects.
type localFilesMonitor struct {}

func (lfm localFilesMonitor) listDir() []os.FileInfo {
	files, err := ioutil.ReadDir("./")
    if err != nil {
        log.Fatal(err)
        files = make([]os.FileInfo, 0)
    }
    return files
}

// Daemon for a periodical check of file changes.
type Checker struct {
	created time.Time
	lastCheckTime time.Time
	checkedFiles map[string]time.Time
	info dirInfo
}

func NewChecker() *Checker {
	return &Checker {
		created: time.Now(),
		checkedFiles: make(map[string]time.Time),
		info: localFilesMonitor{},
	}
}

func (c *Checker) processFiles() {
	files := c.info.listDir()
    c.compare(files)
}

func (c *Checker) compare(files []os.FileInfo) {
	for _, f := range files {
		_, exists := c.checkedFiles[f.Name()]
		if exists && f.ModTime().Unix() > c.lastCheckTime.Unix() {
			fmt.Println("OOOOO:)")
		} else if !exists {
			c.checkedFiles[f.Name()] = f.ModTime()
			fmt.Println("New file found.")
		} else {
			fmt.Println("No new file found.")
		}
	}
}

func (c *Checker) Monitor() {
	for {
		select {
		case <- time.After(2 * time.Second):
			c.processFiles()
			c.lastCheckTime = time.Now()
		}
	}
}

func main() {
	uploadCh := make(chan *models.Command, 1)
	scpCh := make(chan *models.Command, 1)
	go web.RunServer(uploadCh)
	rm := remote.NewRemoteFilesMonitor("root", "root", "0.0.0.0:32768", "/home", scpCh)
	go rm.Run()
	for {
		select {
		case data := <-uploadCh:
			log.Println(*data)
			// Notify remote monitor to copy file on a remote host
			scpCh <- data
		}
	}
	// checker := NewChecker()
	// checker.Monitor()
	// rm := remote.NewRemoteFilesMonitor("root", "root", "0.0.0.0:32768", "/home")
	// rm.Copy("hello.txt", "/home/")	
	
}