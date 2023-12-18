package config

import (
	"fmt"
	"io"
	"os"

	"github.com/olblak/releasepost/internal/core/changelog"
	"gopkg.in/yaml.v3"
)

/*
Config represents the configuration of the application
*/
type Config struct {
	/*
		Changelogs contains the list of changelog to generate.
	*/
	Changelogs []changelog.Config `yaml:"changelogs"`
}

var (
	// ConfigFile is the path to the configuration file
	ConfigFile string
)

// Load loads an releasepost configuration file into memory.
func (c *Config) Load() error {

	f, err := os.Open(ConfigFile)
	if err != nil {
		return fmt.Errorf("opening releasepost configuration file %q: %s", ConfigFile, err)
	}
	defer f.Close()

	configFileByte, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading Updatecli compose file %q: %s", ConfigFile, err)
	}

	err = yaml.Unmarshal(configFileByte, &c)
	if err != nil {
		return fmt.Errorf("parsing Updatecli compose file %q: %s", ConfigFile, err)
	}

	for i := range c.Changelogs {
		err = c.Changelogs[i].Sanitize(ConfigFile)
		if err != nil {
			fmt.Printf("sanitizing changelog configuration for %q: %s", c.Changelogs[i].Name, err)
			continue
		}
	}

	return nil
}
