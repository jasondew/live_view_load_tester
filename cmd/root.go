package cmd

import (
	"os"

	"github.com/jasondew/live_view_load_tester/agent"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "live_view_load_tester",
	Short: "Load test your LiveViews.",
	Long: `Load test your LiveViews.`,
  Run: func(cmd *cobra.Command, args []string) {
		agent.Execute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.live_view_load_tester.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
