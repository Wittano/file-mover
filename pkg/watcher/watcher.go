package watcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/filebot/pkg/config"
	"github.com/wittano/filebot/pkg/path"
	"log"
	"os"
	"path/filepath"
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

func (w *MyWatcher) ObserveFiles() {
	for {
		select {
		case e, ok := <-w.Events:
			if !ok {
				w.blocker <- false
				return
			}

			if e.Has(fsnotify.Create) || e.Has(fsnotify.Rename) {
				if dir, ok := w.fileObserved[e.Name]; ok {
					moveFileToDestination(dir, e.Name)
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

func (w *MyWatcher) AddFilesToObservable(config *config.Config) {
	for _, dir := range config.Dirs {
		for _, src := range dir.Src {
			var (
				err   error
				paths []string
			)

			if dir.Recursive {
				paths, err = path.GetPathFromPatternRecursive(src)
			} else {
				paths, err = path.GetPathsFromPattern(src)
			}

			if err != nil {
				log.Printf("Invalid path: %s", err)
				continue
			}

			if paths != nil {
				go w.fillFileObservedMap(paths, dir.Dest)

				w.addFilesToObservable(paths...)
				go moveFileToDestination(dir.Dest, paths...)
			}
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

func (w *MyWatcher) UpdateObservableFileList(flags config.Flags) {
	var wg sync.WaitGroup

	conf, err := config.Get(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	timer := time.NewTicker(flags.UpdateInterval)
	defer timer.Stop()

	for {
		wg.Add(2)

		<-timer.C

		go w.removeUnnecessaryFiles(&wg)
		go func(wg *sync.WaitGroup) {
			w.AddFilesToObservable(conf)

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

func moveFileToDestination(dest string, paths ...string) {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		log.Printf("Destination directory %s doesn't exist", dest)
		return
	}

	for _, src := range paths {
		_, filename := filepath.Split(src)
		newPath := filepath.Join(dest, filename)

		if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			err := os.Rename(src, newPath)
			if err != nil {
				log.Printf("Failed to move file from %s to %s. %s", src, newPath, err)
				return
			}

			log.Printf("Moved file from %s to %s", src, dest)
		}
	}
}
