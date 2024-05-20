package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pingpong",
	Short: "Pingpong is a simple network latency testing tool",
	Long: `Pingpong is a network latency testing tool written in Go
It can operate in two modes: pinger and ponger.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(pingerCmd)
	rootCmd.AddCommand(pongerCmd)
}
