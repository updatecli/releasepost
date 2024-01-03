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
		Short: "Releasepost is a release note town crier",
		Long: `
Releasepost is a release note town crier.
It retrieves release notes from third location, like GitHub releases,
and then copy them to locally to a directory of your choice.
`,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Releasepost configuration file")
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
