package engine

import (
	"fmt"

	"github.com/updatecli/releasepost/internal/core/runner"
)

/*
Init initializes the releasepost engine.
*/
func (e *Engine) Init() error {

	// Load configuration
	err := e.config.Load()
	if err != nil {
		return fmt.Errorf("loading configuration: %v", err.Error())
	}

	// Init Changelog runners
	for i := range e.config.Changelogs {
		runner, err := runner.New(e.config.Changelogs[i])
		if err != nil {
			fmt.Printf("unable to create runner: %v\n", err.Error())
			continue
		}

		e.runners = append(e.runners, runner)
	}

	if len(e.config.Changelogs) == 0 {
		return fmt.Errorf("no changelog found in configuration file")
	}

	return nil
}
