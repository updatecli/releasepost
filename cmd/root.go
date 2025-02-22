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
	configFile       string
	e                engine.Engine
	dryRun           bool
	cleanRun         bool
	detailedExitCode bool

	detailedExitCodeCmdDescription string = `Returns a detailed exit code when the command exits.
When provided, this argument changes the exit codes and their meanings
to provide more granular information about what the resulting plan contains:

0 = Succeeded with empty diff (no changes)
1 = Error
2 = Succeeded with non-empty diff (changes present)
`

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
			err = e.Run(cleanRun)
			if cleanRun {
				fmt.Println("Clean run mode enabled, releasepost will remove any files not created by releasepost in changelogs directories !")
			}
			if err != nil {
				fmt.Printf("Failed to run releasepost: %v", err)
				os.Exit(1)
			}

			os.Exit(result.ChangelogResult.ExitCode(detailedExitCode))
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
	rootCmd.PersistentFlags().BoolVar(&cleanRun, "clean", false, "Clean run, removes files from changelog directories not created by releasepost.")
	rootCmd.PersistentFlags().BoolVar(&detailedExitCode, "detailed-exit-code", false, detailedExitCodeCmdDescription)
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
