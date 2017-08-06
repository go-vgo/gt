// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

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
