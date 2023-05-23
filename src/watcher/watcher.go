package watcher

import (
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/path"
	"log"

	"github.com/fsnotify/fsnotify"
)

type MyWatcher struct {
	*fsnotify.Watcher
	blocker chan bool
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

			log.Printf("Event %s", e)
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
				w.addFilesToObservable(paths...)
			}
		}
	}
}

func (w MyWatcher) addFilesToObservable(paths ...string) {
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			log.Printf("Cannot add %s file/directory to tracing list: %s", p, err)
		}
	}
}
