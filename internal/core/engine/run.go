package engine

import (
	"fmt"

	"github.com/updatecli/releasepost/internal/core/result"
)

/*
Run executes the engine.
It will run all the runners and save the changelogs to disk.
*/
func (e *Engine) Run(cleanRun bool) error {

	for i := range e.config.Changelogs {
		changelogs, err := e.runners[i].Run()
		if err != nil {
			fmt.Printf("unable to run runner: %v\n", err.Error())
			continue
		}

		if len(changelogs) == 0 {
			fmt.Printf("no changelog found for %s\n", e.config.Changelogs[i].Name)
			continue
		}

		err = e.config.Changelogs[i].SaveToDisk(changelogs)
		if err != nil {
			fmt.Printf("unable to save changelog to disk: %v\n", err.Error())
			continue
		}

		err = e.config.Changelogs[i].SaveIndexToDisk(changelogs)
		if err != nil {
			fmt.Printf("unable to save changelog index to disk: %v\n", err.Error())
			continue
		}

		err = result.ChangelogResult.UpdateResult(&e.result)
		if err != nil {
			fmt.Printf("unable to update result: %v\n", err.Error())
			continue
		}
	}

	if err := e.Clean(cleanRun); err != nil {
		fmt.Printf("unable to clean changelog: %v\n", err.Error())
	}

	fmt.Println(e.result)

	return nil
}
