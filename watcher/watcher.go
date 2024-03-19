package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/wittano/filebot/file"
	"github.com/wittano/filebot/setting"
	"github.com/wittano/filebot/tasks"
	"os"
	"sync"
)

type MyWatcher struct {
	*fsnotify.Watcher
	ctx          context.Context
	mutex        sync.Mutex
	blocker      chan bool
	fileObserved map[string]string
}

func (w *MyWatcher) Close() (err error) {
	w.mutex.Lock()
	close(w.blocker)

	err = w.Watcher.Close()
	w.mutex.Unlock()

	return
}

func NewWatcher(ctx context.Context) MyWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		setting.Logger().Fatal("Failed initialized system file watcher", err)
	}

	blocker := make(chan bool)

	return MyWatcher{
		w,
		ctx,
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
					if err := file.MoveToDestination(dir, e.Name); err != nil {
						setting.Logger().Error(fmt.Sprintf("Failed move file from %s to %s", dir, e.Name), err)
					}
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

func (w *MyWatcher) WaitForEvents() {
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

			go func(dest string, srcs []string) {
				if err = file.MoveToDestination(dest, srcs...); err != nil {
					setting.Logger().Error(fmt.Sprintf("One of soruce file wasn't moved to destination directory"), err)
					return
				}
			}(destPath, paths)
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
			setting.Logger().Error(fmt.Sprintf("Cannot add %s file/directory to tracing list", p), err)
		}
	}
}

func (w *MyWatcher) UpdateObservableFileList() {
	tasks.RunTaskWithInterval(w.ctx, setting.Flags.UpdateInterval, w.updateObservableFileList)
}

func (w *MyWatcher) updateObservableFileList(ctx context.Context) error {
	select {
	case <-w.ctx.Done():
		return nil
	default:
	}

	var wg sync.WaitGroup
	wg.Add(2) // Add number of task, that we should end before we can continue task

	go w.removeUnnecessaryFiles(ctx, &wg)
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			conf := setting.Flags.Config()

			w.AddFilesToObservable(*conf)
		}
	}()

	wg.Wait()

	return nil
}

func (w *MyWatcher) removeUnnecessaryFiles(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
		for _, path := range w.WatchList() {
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				w.Remove(path)
			}
		}
	}
}
