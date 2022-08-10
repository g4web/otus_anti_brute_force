package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/proto"
	"github.com/spf13/cobra"
)

var removeNetworkFromBlacklistCmd = &cobra.Command{
	Use:   "removeNetworkFromBlacklist",
	Short: "Remove a network from the blacklist",
	Long:  `Remove a network from the blacklist`,
	Run: func(cmd *cobra.Command, args []string) {
		if network == "" {
			log.Fatal(errors.New("network is not specified"))
			return
		}

		r, err := client.RemoveNetworkFromBlackList(ctx, &proto.RemoveNetworkFromBlackListRequest{Network: network})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeNetworkFromBlacklistCmd)

	removeNetworkFromBlacklistCmd.PersistentFlags().StringVar(
		&network,
		"n",
		"",
		"Network, for example \"192.168.0.0/24\"",
	)
}
