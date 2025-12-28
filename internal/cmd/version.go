package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information - set during build
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print version, commit, and build information for dsaprep.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dsaprep %s\n", Version)
		fmt.Printf("  Commit:     %s\n", Commit)
		fmt.Printf("  Built:      %s\n", BuildDate)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
