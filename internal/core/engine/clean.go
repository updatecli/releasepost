package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/updatecli/releasepost/internal/core/result"
)

// Clean analyze changelog monitored directories and remove any changelog that
// weren't created or modified by releasepost.
func (e *Engine) Clean(clean bool) error {

	fmt.Printf("\n\nCleaning\n")
	for dirpath, dir := range e.result.Dir {
		files, err := filepath.Glob(filepath.Join(dirpath, "*"))
		if err != nil {
			return err
		}

		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				fmt.Printf("unable to get file info: %v\n", err.Error())
				continue
			}

			if info.IsDir() {
				continue
			}

			if !dir.IsFileExist(f) {
				err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if info.IsDir() {
						return nil
					}

					result.ChangelogResult.UnTracked = append(result.ChangelogResult.UnTracked, f)
					if clean {
						fmt.Printf("\t* removing %s\n", f)
						err = os.Remove(path)
						if err != nil {
							return err
						}
					}

					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	if len(result.ChangelogResult.UnTracked) > 0 {
		fmt.Printf("Untracked files detected:\n")
		for _, f := range result.ChangelogResult.UnTracked {
			fmt.Printf("\t* %s\n", f)
		}
	}

	fmt.Printf("\n\n")
	return nil
}
