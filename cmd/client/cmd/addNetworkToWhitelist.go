package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/proto"
	"github.com/spf13/cobra"
)

var addNetworkToWhitelistCmd = &cobra.Command{
	Use:   "addNetworkToWhitelist",
	Short: "Add a network to the whitelist",
	Long:  `Add a network to the whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		if network == "" {
			log.Fatal(errors.New("network is not specified"))
			return
		}

		r, err := client.AddNetworkToWhiteList(ctx, &proto.AddNetworkToWhiteListRequest{Network: network})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(addNetworkToWhitelistCmd)

	addNetworkToWhitelistCmd.PersistentFlags().StringVar(
		&network,
		"n",
		"",
		"Network, for example \"192.168.0.0/24\"",
	)
}
