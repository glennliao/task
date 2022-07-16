package glob

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"
)

func matchPattern(dir string, pattern *ParsedPattern, ignore Patterns, absolutePaths bool) ([]string, error) {
	inputPath := filepath.Join(dir, pattern.Input)
	if _, err := os.Stat(inputPath); err == nil {

		if absolutePaths {
			return []string{inputPath}, nil
		}
		return []string{filepath.Join(pattern.Input)}, nil
	}

	startDir := dir
	if len(pattern.Base) > 0 {
		startDir = filepath.Join(dir, pattern.Base)
	}

	var files []string
	err := filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		for _, igr := range ignore {
			if Match(igr, rel) {

				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		if !d.IsDir() {
			if Match(pattern, rel) {
				if absolutePaths {
					files = append(files, path)
				} else {
					files = append(files, rel)
				}
			}

		}
		return nil
	})
	return files, err
}

func matchFiles(dir string, globs Patterns, ignore Patterns, absolutePaths bool) ([]string, error) {
	var files []string
	var mu sync.Mutex
	var g errgroup.Group

	for _, pattern := range globs {
		pattern := pattern
		g.Go(func() error {
			matched, err := matchPattern(dir, pattern, ignore, absolutePaths)
			mu.Lock()
			files = append(files, matched...)
			mu.Unlock()

			return err
		})
	}

	return files, g.Wait()
}

func Match(p *ParsedPattern, path string) bool {
	return p.RegExp.MatchString(path)
}
