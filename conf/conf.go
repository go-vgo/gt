package conf

import (
	"fmt"
	"log"
	"sync"
	// "time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

// NewWatcher new fsnotify watcher
func NewWatcher(paths string, v interface{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				// if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				// 	log.Println("watcher.Events: ignore CHMOD event:", event)
				// 	continue
				// }

				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("modified file:", event.Name)
					Init(paths, v)
					// log.Println("watch-config...", config)
					log.Println("watch-config...", v)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(paths)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

var (
	// config     Config
	configLock = new(sync.RWMutex)
)

// Init toml config
func Init(tpath string, v interface{}) {
	configLock.Lock()
	if _, err := toml.DecodeFile(tpath, v); err != nil {
		fmt.Println(err)
		return
	}
	configLock.Unlock()
}
