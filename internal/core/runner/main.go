package runner

import (
	"fmt"

	"github.com/olblak/releasepost/internal/core/changelog"
	"github.com/olblak/releasepost/internal/plugins/github"
)

/*
Runner is the interface that wraps the Run method.
*/
type Runner interface {
	Run() ([]changelog.Spec, error)
}

/*
New returns a new Runner based on the configuration.
*/
func New(s changelog.Config) (Runner, error) {
	switch s.Kind {
	case "github":
		return github.New(s.Spec)
	default:
		return nil, fmt.Errorf("unknown changelog kind %q", s.Kind)
	}
}
