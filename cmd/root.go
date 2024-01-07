package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/updatecli/releasepost/internal/core/config"
	"github.com/updatecli/releasepost/internal/core/dryrun"
	"github.com/updatecli/releasepost/internal/core/engine"
	"github.com/updatecli/releasepost/internal/core/result"
)

var (
	configFile string
	e          engine.Engine
	dryRun     bool

	rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {

			if configFile == "" {
				fmt.Println("missing configuration file")
				os.Exit(2)
			}

			config.ConfigFile = configFile
			err := e.Init()
			if err != nil {
				fmt.Printf("Failed to initialize releasepost: %v", err)
				os.Exit(2)
			}

			fmt.Println("Running releasepost")
			if dryRun {
				dryrun.Enabled = true
				fmt.Println("Dry run mode enabled, no changelog will be saved to disk")
			}
			err = e.Run()
			if err != nil {
				fmt.Printf("Failed to run releasepost: %v", err)
				os.Exit(2)
			}

			os.Exit(result.ChangelogResult.ExitCode())
		},
		Use:   "releasepost",
		Short: "releasepost is a release note town crier",
		Long: `
releasepost is a release note town crier.
It retrieves release notes from a third location, like a GitHub release,
and then copy them locally to your directory of choice.
It creates one file per release note version and an index file.
It can creates files using different formats like markdown, asciidoctor, or json.
`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "releasepost configuration file")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	rootCmd.AddCommand(
		versionCmd,
	)
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to execute releasepost: %v", err)
		os.Exit(1)
	}
}
