package result

import (
	"fmt"
)

type FileResult struct {
	State string
	Path  string
}

type FileResults []FileResult

type Result struct {
	Dir map[string]FileResults
}

func (f FileResults) IsFileExist(file string) bool {
	for _, f := range f {
		if f.Path == file {
			return true
		}
	}
	return false
}

func (f Result) String() string {
	s := "Result\n"
	for _, files := range f.Dir {
		for _, f := range files {
			s = fmt.Sprintf("%s\t* %s (%s)\n", s, f.Path, f.State)
		}
	}
	return s
}
