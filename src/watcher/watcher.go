package watcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/path"
	"log"
	"os"
	p "path"
)

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
				moveFileToDestination(w.fileObserved[e.Name], e.Name)
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

func (w MyWatcher) AddFilesToObservable(config config.Config) {
	for _, dir := range config.Dirs {
		for _, src := range dir.Src {
			paths, err := path.GetPathsFromPattern(src)
			if err != nil {
				log.Fatalf("Invalid path: %s", err)
			}

			if paths != nil {
				go w.fillFileObserved(paths, dir.Dest)

				w.addFilesToObservable(paths...)
				go moveFileToDestination(dir.Dest, paths...)
			}
		}
	}
}

func (w MyWatcher) fillFileObserved(src []string, dest string) {
	for _, path := range src {
		w.fileObserved[path] = dest
	}
}

func (w MyWatcher) addFilesToObservable(paths ...string) {
	for _, path := range paths {
		if err := w.Add(path); err != nil {
			log.Printf("Cannot add %s file/directory to tracing list: %s", path, err)
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
			log.Printf("Moved file from %s to %s", newPath, dest)
		}
	}
}
