package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Ehco1996/clash-speed/internal/ui"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:               "clash-speed",
	Short:             "use clash-speed to test your proxy server speed over www.speedtest.net",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {

		if cfgFile == "" {
			fmt.Println("cfg file path is empty")
			os.Exit(1)
		}
		if err := runTUI(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"clash config file path (also support download from http such as your clash subscribe link)")

	rootCmd.PersistentFlags().IntVar(&ui.DownLoadConcurrency, "concurrency", 1,
		"tcp concurrency, note that more concurrency will use more traffic")

	rootCmd.PersistentFlags().DurationVar(&ui.TestDuration, "duration", 3*time.Second,
		"time duration for speed test")
}
