package cmd

import (
	"pingpong/pingpong"

	"github.com/spf13/cobra"
)

var pingerCmd = &cobra.Command{
	Use:   "pinger",
	Short: "Pinger sends packets to a target and measures latency",
	Long: `Pinger sends packets containing timestamps to a specified target,
which should be running in ponger mode. It measures and displays the latency percentiles.`,
	Run: func(cmd *cobra.Command, args []string) {
		pingpong.RunPinger()
	},
}

func init() {
	pingerCmd.Flags().StringP("target", "t", "", "Target IP:Port to send packets to")
	pingerCmd.Flags().IntP("numPackets", "n", 100, "Number of packets to send (default: 100)")
	pingerCmd.Flags().IntP("numPorts", "p", 50, "Number of random source ports (default: 50)")
	pingerCmd.Flags().StringP("protocol", "x", "udp", "Protocol to use (udp or tcp, default: udp)")

	pingerCmd.MarkFlagRequired("target")
}
