package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"github.com/spf13/cobra"
)

var removeNetworkFromWhitelistCmd = &cobra.Command{
	Use:   "removeNetworkFromWhitelist",
	Short: "Remove a network from the whitelist",
	Long:  `Remove a network from the whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		if network == "" {
			log.Fatal(errors.New("network is not specified"))
			return
		}

		r, err := client.RemoveNetworkFromWhiteList(ctx, &proto.RemoveNetworkFromWhiteListRequest{Network: network})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeNetworkFromWhitelistCmd)

	removeNetworkFromWhitelistCmd.PersistentFlags().StringVar(
		&network,
		"n",
		"",
		"Network, for example \"192.168.0.0/24\"",
	)
}
