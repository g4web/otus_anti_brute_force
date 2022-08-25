package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"github.com/spf13/cobra"
)

var addNetworkToBlacklistCmd = &cobra.Command{
	Use:   "addNetworkToBlacklist",
	Short: "Add a network to the whitelist",
	Long:  `Add a network to the whitelist`,
	Run: func(cmd *cobra.Command, args []string) {
		if network == "" {
			log.Fatal(errors.New("network is not specified"))
			return
		}

		r, err := client.AddNetworkToBlackList(ctx, &proto.AddNetworkToBlackListRequest{Network: network})

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(addNetworkToBlacklistCmd)

	addNetworkToBlacklistCmd.PersistentFlags().StringVar(
		&network,
		"n",
		"",
		"Network, for example \"192.168.0.0/24\"",
	)

	addNetworkToBlacklistCmd.PersistentFlags().StringVar(
		&configFile,
		"c",
		"../../configs/config.env",
		"A path to config file",
	)
}
