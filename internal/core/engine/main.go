package engine

import (
	"github.com/olblak/releasepost/internal/core/config"
	"github.com/olblak/releasepost/internal/core/runner"
)

/*
Engine is the main structure of the application
*/
type Engine struct {
	config  config.Config
	runners []runner.Runner
}
