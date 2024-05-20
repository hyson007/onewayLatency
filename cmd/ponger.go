package cmd

import (
	"pingpong/pingpong"

	"github.com/spf13/cobra"
)

var pongerCmd = &cobra.Command{
	Use:   "ponger",
	Short: "Ponger receives packets from a pinger and sends them back",
	Long: `Ponger listens for incoming packets from a pinger, appends its own timestamp,
and sends the packet back.`,
	Run: func(cmd *cobra.Command, args []string) {
		pingpong.RunPonger()
	},
}

func init() {
	pongerCmd.Flags().IntP("port", "p", 0, "Listening port for incoming packets")
	pongerCmd.Flags().StringP("protocol", "x", "udp", "Protocol to use (udp or tcp, default: udp)")

	pongerCmd.MarkFlagRequired("port")
}
