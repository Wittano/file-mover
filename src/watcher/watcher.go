package watcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/path"
	"log"
	"os"
	p "path"
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

func (w *MyWatcher) AddFilesToObservable(config config.Config) {
	for _, dir := range config.Dirs {
		for _, src := range dir.Src {
			var paths []string
			var err error

			if dir.Recursive {
				paths, err = path.GetPathFromPatternRecursive(src)
			} else {
				paths, err = path.GetPathsFromPattern(src)
			}

			if err != nil {
				log.Fatalf("Invalid path: %s", err)
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
	for _, path := range src {
		mutex.Lock()
		w.fileObserved[path] = dest
		mutex.Unlock()
	}
}

func (w *MyWatcher) addFilesToObservable(paths ...string) {
	for _, path := range paths {
		if err := w.Add(path); err != nil {
			log.Printf("Cannot add %s file/directory to tracing list: %s", path, err)
		}
	}
}

func (w *MyWatcher) UpdateObservableFileList(flags config.FlagConfig) {
	var wg sync.WaitGroup

	for {
		wg.Add(2)

		time.Sleep(flags.UpdateInterval)

		go w.removeUnnecessaryFiles(&wg)
		go func(wg *sync.WaitGroup) {
			conf, err := config.LoadConfig(flags.ConfigPath)
			if err != nil {
				log.Fatalf("Failed loaded configuration: %s", err)
			}

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
		log.Printf("Destination directory %s not exist", dest)
		return
	}

	for _, src := range paths {
		_, filename := p.Split(src)
		newPath := dest + "/" + filename

		if _, err := os.Stat(src); !errors.Is(err, os.ErrNotExist) {
			os.Rename(src, newPath)
			log.Printf("Moved file from %s to %s", src, dest)
		}
	}
}
