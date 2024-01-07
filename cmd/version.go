package cmd

import (
	"strings"

	"fmt"

	"github.com/spf13/cobra"

	"github.com/updatecli/releasepost/internal/core/version"
)

var (
	// Version Contains application version
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print current application version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n%s\n", strings.ToTitle("Version"))
			version.Show()
		},
	}
)
