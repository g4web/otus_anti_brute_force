package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"github.com/spf13/cobra"
)

var ip string

var deleteIPStatsCmd = &cobra.Command{
	Use:   "deleteIpStats",
	Short: "delete statistics for ip",
	Run: func(cmd *cobra.Command, args []string) {
		if ip == "" {
			log.Fatal(errors.New("IP is not specified"))
			return
		}

		r, err := client.DeleteIPStats(ctx, &proto.DeleteIPStatsRequest{IP: ip})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteIPStatsCmd)

	deleteIPStatsCmd.PersistentFlags().StringVar(
		&ip,
		"ip",
		"",
		"IP address, for example \"192.168.0.1\"",
	)
}
