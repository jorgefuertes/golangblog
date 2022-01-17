package dev

import (
	"time"

	"golangblog/internal/log"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"github.com/karrick/godirwalk"
)

var lr *lrserver.Server

// DelayedReload does a reload after the sec seconds
func DelayedReload(sec time.Duration, reason string) {
	go func() {
		time.Sleep(sec * time.Second)
		lr.Reload(reason)
	}()
}

// StartLiveReload starts lr-server and watchers
func StartLiveReload() {
	// Watched folders
	watchedDirs := []string{"static", "view"}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("LR", err)
		return
	}

	// Create and start LiveReload server
	lr = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	lr.SetLiveCSS(true)
	go lr.ListenAndServe()

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debug("LR", event.Name)
					DelayedReload(1, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("LR", err)
			}
		}
	}()

	for _, v := range watchedDirs {
		// Watch subdirs
		godirwalk.Walk(v, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if de.IsDir() {
					log.Debug("LR", osPathname)
					err = watcher.Add(osPathname)
					if err != nil {
						log.Error("LR", err)
					}
				}
				return nil
			},
			Unsorted: true,
		})
	}
	// Sending a first reload after initialization
	DelayedReload(2, "System restart")
}
