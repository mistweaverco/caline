package caline

import (
	"os"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/caline/internal/config"
	"github.com/mistweaverco/caline/internal/overlay"
	"github.com/spf13/cobra"
)

var VERSION string
var cfg = config.NewConfig(config.Config{})

var rootCmd = &cobra.Command{
	Use:   "caline",
	Short: "A simple tool to display calendar events in an overlay line",
	Long:  "A simple tool to display calendar events in an overlay line",
	Run: func(cmd *cobra.Command, files []string) {
		if cfg.Flags.Version {
			log.Info("Version", runtime.GOOS, VERSION)
			return
		}
		overlay.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&cfg.Flags.Version, "version", false, "version")
}
