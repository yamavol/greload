package lib

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/yamavol/greload/log"
)

type WatchOptions struct {
	Dirs   []string
	Ignore []string
}

func listSubdir(dirs ...string) ([]string, error) {

	dirList := make([]string, 0)
	dirSet := make(map[string]bool)

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && !dirSet[path] {
				dirList = append(dirList, path)
				dirSet[path] = true
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return dirList, nil
}

func excludeSubdir(dirs []string, exclude []string) []string {
	excludeSet := make(map[string]bool)
	for _, dir := range exclude {
		excludeSet[dir] = true
	}

	filtered := make([]string, 0)
	for _, dir := range dirs {
		if !excludeSet[dir] {
			filtered = append(filtered, dir)
		}
	}
	return filtered
}

// Start filesystem watcher.
func WatchStart(watch []string, exclude []string, srv *ProxyServer) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Error(ok)
					return
				}
				log.Debug("[fs]", "event", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					srv.TriggerReload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("[fs]", "error", err)
			}
		}
	}()

	dirList, err := listSubdir(watch...)
	if err != nil {
		log.Error(err)
		return
	}
	excludeList, err := listSubdir(exclude...)
	if err != nil {
		log.Error(err)
	}
	watchList := excludeSubdir(dirList, excludeList)

	if err != nil {
		log.Error(err)
	}

	for _, dir := range watchList {
		log.Debug("[fs]", "watch", dir)
		err = watcher.Add(dir)
	}
	if err != nil {
		log.Error(err)
	}
	<-done
}
