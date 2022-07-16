package glob

import (
	"github.com/pkg/errors"
	"os"
)

const (
	STAR     = `([^/]*)`
	GLOBSTER = `((?:[^/]*(?:\/|$))*)`
)

type Patterns = map[string]*ParsedPattern

type Options struct {
	IgnorePatterns []string
	CWD            string
	Patterns       []string
	AbsolutePaths  bool
}

func Pattern(pattern string) *Options {
	return &Options{
		Patterns: []string{pattern},
	}
}

func CWD(cwd string) *Options {
	return &Options{
		CWD: cwd,
	}
}

func IgnorePattern(pattern string) *Options {
	return &Options{
		IgnorePatterns: []string{pattern},
	}
}

func Glob(options ...*Options) ([]string, error) {
	var cwd string
	ignores := make(Patterns)
	patterns := make(Patterns)
	absolutePaths := false

	for _, opt := range options {
		if len(opt.CWD) > 0 {
			cwd = opt.CWD
		}

		if opt.AbsolutePaths {
			absolutePaths = true
		}

		for _, p := range opt.IgnorePatterns {
			reg, err := Parse(p)
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse ignorePattern %s", p)
			}
			ignores[p] = reg
		}

		for _, p := range opt.Patterns {
			reg, err := Parse(p)

			if err != nil {
				return nil, errors.Wrapf(err, "could not parse provided pattern %s", p)
			}

			patterns[p] = reg
		}
	}

	if len(patterns) < 1 {
		return nil, errors.New("No patterns provided! Please provide a valid glob pattern as parameter")
	}

	if len(cwd) < 1 {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "could not determine current working directory, please provide it as argument")
		}
		cwd = wd
	}

	return matchFiles(cwd, patterns, ignores, absolutePaths)
}
