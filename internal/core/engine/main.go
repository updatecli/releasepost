package engine

import (
	"github.com/updatecli/releasepost/internal/core/config"
	"github.com/updatecli/releasepost/internal/core/result"
	"github.com/updatecli/releasepost/internal/core/runner"
)

/*
Engine is the main structure of the application
*/
type Engine struct {
	config  config.Config
	runners []runner.Runner
	result  result.Result
}
