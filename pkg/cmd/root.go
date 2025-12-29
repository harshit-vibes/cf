// Package cmd provides CLI commands for dsaprep
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/harshit-vibes/dsaprep/pkg/external/cfapi"
	"github.com/harshit-vibes/dsaprep/pkg/external/cfweb"
	exthealth "github.com/harshit-vibes/dsaprep/pkg/external/health"
	"github.com/harshit-vibes/dsaprep/pkg/internal/config"
	"github.com/harshit-vibes/dsaprep/pkg/internal/health"
	"github.com/harshit-vibes/dsaprep/pkg/internal/workspace"
)

var (
	// Version information
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"

	// Command line flags
	skipChecks bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "dsaprep",
	Short: "DSA Practice CLI - Your competitive programming companion",
	Long: `dsaprep is a CLI tool for practicing Data Structures and Algorithms
with Codeforces integration. It provides:

  • Problem parsing and workspace management
  • Solution submission and verdict tracking
  • Practice progress tracking
  • Beautiful TUI dashboard`,
	PersistentPreRunE: runPreChecks,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Launch TUI
		fmt.Println("DSA Prep CLI v" + Version)
		fmt.Println("TUI not yet implemented. Use --help for available commands.")
		return nil
	},
}

func runPreChecks(cmd *cobra.Command, args []string) error {
	// Skip health checks for version and help commands
	if cmd.Name() == "version" || cmd.Name() == "help" {
		return nil
	}
	return runStartupChecks()
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Add flags
	rootCmd.PersistentFlags().BoolVar(&skipChecks, "skip-checks", false, "Skip startup health checks")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(parseCmd)
}

func initConfig() {
	// Load configuration
	if err := config.Init(""); err != nil {
		// Config init failure is handled by health checks
		return
	}
}

func runStartupChecks() error {
	if skipChecks {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	checker := health.NewChecker()

	// Get workspace path
	cfg := config.Get()
	wsPath := "."
	if cfg != nil && cfg.WorkspacePath != "" {
		wsPath = cfg.WorkspacePath
	}
	ws := workspace.New(wsPath)

	// Internal checks
	checker.AddCheck(&health.EnvFileCheck{})
	checker.AddCheck(&health.ConfigCheck{})
	checker.AddCheck(health.NewWorkspaceCheck(ws))
	checker.AddCheck(health.NewSchemaVersionCheck(ws))
	checker.AddCheck(&health.SessionCheck{})

	// External checks (only if we have credentials)
	creds, _ := config.LoadCredentials()
	var apiClient *cfapi.Client
	if creds != nil && creds.IsAPIConfigured() {
		apiClient = cfapi.NewClient(cfapi.WithAPICredentials(creds.APIKey, creds.APISecret))
	} else {
		apiClient = cfapi.NewClient()
	}

	parser := cfweb.NewParserWithClient(nil)

	checker.AddCheck(exthealth.NewCFAPICheck(apiClient))
	checker.AddCheck(exthealth.NewCFWebCheck(parser))
	checker.AddCheck(exthealth.NewCFHandleCheck(apiClient))

	// Run checks
	report := checker.Run(ctx)

	// Display results
	if verbose || report.OverallStatus != health.StatusHealthy {
		displayHealthReport(report)
	}

	if !report.CanProceed {
		fmt.Println("\n❌ Cannot proceed due to critical errors. Please fix the issues above.")
		return fmt.Errorf("startup checks failed")
	}

	if report.OverallStatus == health.StatusDegraded {
		fmt.Println("\n⚠️  Some features may be unavailable. See warnings above.")
	}

	return nil
}

func displayHealthReport(report *health.Report) {
	fmt.Printf("\n🔍 Health Check Report (took %s)\n", report.Duration.Round(time.Millisecond))
	fmt.Println("─────────────────────────────────")

	for _, result := range report.Results {
		var icon string
		switch result.Status {
		case health.StatusHealthy:
			icon = "✓"
		case health.StatusDegraded:
			icon = "⚠"
		case health.StatusCritical:
			icon = "✗"
		}

		fmt.Printf("%s %-20s %s\n", icon, result.Name, result.Message)
		if result.Details != "" && result.Status != health.StatusHealthy {
			fmt.Printf("  └─ %s\n", result.Details)
		}
	}

	fmt.Println("─────────────────────────────────")
	fmt.Printf("Status: %s | Schema: %s\n", report.OverallStatus, report.CurrentSchemaVersion)
}

// versionCmd shows version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dsaprep %s\n", Version)
		fmt.Printf("  Commit:     %s\n", Commit)
		fmt.Printf("  Build Date: %s\n", BuildDate)
	},
}

// initCmd initializes the workspace
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new dsaprep workspace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		ws := workspace.New(path)

		if ws.Exists() {
			return fmt.Errorf("workspace already exists at %s", path)
		}

		creds, _ := config.LoadCredentials()
		handle := ""
		if creds != nil {
			handle = creds.CFHandle
		}

		if err := ws.Init("DSA Practice", handle); err != nil {
			return fmt.Errorf("failed to initialize workspace: %w", err)
		}

		fmt.Printf("✓ Initialized workspace at %s\n", path)
		return nil
	},
}

// healthCmd shows health status
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check system health",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip the regular pre-run checks for health command
		// since we'll run them ourselves with verbose=true
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Force verbose output
		verbose = true
		return runStartupChecks()
	},
}

// parseCmd parses a problem from CF
var parseCmd = &cobra.Command{
	Use:   "parse <contest_id> <problem_index>",
	Short: "Parse a problem from Codeforces",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var contestID int
		if _, err := fmt.Sscanf(args[0], "%d", &contestID); err != nil {
			return fmt.Errorf("invalid contest ID: %s", args[0])
		}
		problemIndex := args[1]

		parser := cfweb.NewParserWithClient(nil)
		problem, err := parser.ParseProblem(contestID, problemIndex)
		if err != nil {
			return fmt.Errorf("failed to parse problem: %w", err)
		}

		fmt.Printf("✓ Parsed: %s. %s\n", problem.Index, problem.Name)
		fmt.Printf("  Rating: %d | Time: %s | Memory: %s\n",
			problem.Rating, problem.TimeLimit, problem.MemoryLimit)
		fmt.Printf("  Tags: %v\n", problem.Tags)
		fmt.Printf("  Samples: %d\n", len(problem.Samples))

		// Save to workspace if available
		cfg := config.Get()
		if cfg != nil && cfg.WorkspacePath != "" {
			ws := workspace.New(cfg.WorkspacePath)
			if ws.Exists() {
				schemaProblem := problem.ToSchemaProblem()
				if err := ws.SaveProblem(schemaProblem); err != nil {
					return fmt.Errorf("failed to save problem: %w", err)
				}
				fmt.Printf("✓ Saved to workspace\n")
			}
		}

		return nil
	},
}
