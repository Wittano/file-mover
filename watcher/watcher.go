package watcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/filebot/cron"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"log"
	"os"
	"sync"
	"time"
)

var mutex sync.Mutex

type MyWatcher struct {
	*fsnotify.Watcher
	blocker      chan bool
	fileObserved map[string]string
}

func NewWatcher() MyWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed initialized file system w: %s", err)
	}

	blocker := make(chan bool)

	return MyWatcher{
		w,
		blocker,
		make(map[string]string),
	}
}

func (w MyWatcher) ObserveFiles() {
	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				w.blocker <- false
				return
			}

			if e.Has(fsnotify.Create) || e.Has(fsnotify.Rename) {
				if dir, ok := w.fileObserved[e.Name]; ok {
					file.MoveToDestination(dir, e.Name)
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				w.blocker <- false
				return
			}

			log.Printf("Error %s", err)
		}
	}
}

func (w MyWatcher) WaitForEvents() {
	if ok := <-w.blocker; !ok {
		return
	}
}

func (w *MyWatcher) AddFilesToObservable(config setting.Config) {
	for _, dir := range config.Dirs {
		paths, err := dir.RealPaths()
		if err != nil {
			log.Printf("Failed to get path for files")
			continue
		}

		if paths != nil {
			destPath := dir.Dest
			if dir.Dest == "" {
				destPath = cron.TrashPath
			}

			go w.fillFileObservedMap(paths, destPath)

			w.addFilesToObservable(paths...)
			go file.MoveToDestination(destPath, paths...)
		}
	}
}

func (w *MyWatcher) fillFileObservedMap(src []string, dest string) {
	for _, p := range src {
		mutex.Lock()
		w.fileObserved[p] = dest
		mutex.Unlock()
	}
}

func (w *MyWatcher) addFilesToObservable(paths ...string) {
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			log.Printf("Cannot add %s file/directory to tracing list: %s", p, err)
		}
	}
}

// TODO Migrate schedule task to go-cron task
func (w *MyWatcher) UpdateObservableFileList() {
	var wg sync.WaitGroup

	conf := setting.Flags.Config()

	timer := time.NewTicker(setting.Flags.UpdateInterval)
	defer timer.Stop()

	for {
		wg.Add(2)

		<-timer.C

		go w.removeUnnecessaryFiles(&wg)
		go func(wg *sync.WaitGroup) {
			w.AddFilesToObservable(*conf)

			wg.Done()
		}(&wg)

		wg.Wait()
	}
}

func (w *MyWatcher) removeUnnecessaryFiles(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, path := range w.WatchList() {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			w.Remove(path)
		}
	}
}
