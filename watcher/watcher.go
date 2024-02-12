package watcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"os"
	"sync"
	"time"
)

type MyWatcher struct {
	*fsnotify.Watcher
	mutex        sync.Mutex
	blocker      chan bool
	fileObserved map[string]string
}

func (w *MyWatcher) Close() (err error) {
	w.mutex.Lock()
	close(w.blocker)

	err = w.Watcher.Close()
	w.mutex.Lock()

	return
}

func NewWatcher() MyWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		setting.Logger().Fatal("Failed initialized system file watcher", err)
	}

	blocker := make(chan bool)

	return MyWatcher{
		w,
		sync.Mutex{},
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

			setting.Logger().Error("Watcher got unexpected error", err)
		}
	}
}

func (w MyWatcher) WaitForEvents() {
	if ok := <-w.blocker; !ok {
		w.Close()

		return
	}
}

func (w *MyWatcher) AddFilesToObservable(config setting.Config) {
	for _, dir := range config.Dirs {
		paths, err := dir.RealPaths()
		if err != nil {
			setting.Logger().Error("Failed to get path for files", err)
			continue
		}

		if paths != nil {
			destPath := dir.Dest
			if dir.Dest == "" {
				destPath, err = dir.TrashDir()
			}

			if err != nil {
				setting.Logger().Error("Failed to find trash directory", err)
				break
			}

			go w.fillFileObservedMap(paths, destPath)

			w.addFilesToObservable(paths...)
			go file.MoveToDestination(destPath, paths...)
		}
	}
}

func (w *MyWatcher) fillFileObservedMap(src []string, dest string) {
	for _, p := range src {
		w.mutex.Lock()
		w.fileObserved[p] = dest
		w.mutex.Unlock()
	}
}

func (w *MyWatcher) addFilesToObservable(paths ...string) {
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			setting.Logger().Errorf("Cannot add %s file/directory to tracing list", err, p)
		}
	}
}

func (w *MyWatcher) UpdateObservableFileList() {
	var wg sync.WaitGroup

	conf := setting.Flags.Config()

	timer := time.NewTicker(setting.Flags.UpdateInterval)
	defer timer.Stop()

	ok := true

	for ok {
		wg.Add(2)

		<-timer.C

		go w.removeUnnecessaryFiles(&wg)
		go func(wg *sync.WaitGroup) {
			w.AddFilesToObservable(*conf)

			wg.Done()
		}(&wg)

		wg.Wait()
		_, ok = <-w.blocker
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
