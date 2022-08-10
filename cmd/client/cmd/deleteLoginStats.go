package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/g4web/otus_anti_brute_force/proto"
	"github.com/spf13/cobra"
)

var login string

var deleteLoginStatsCmd = &cobra.Command{
	Use:   "deleteLoginStats",
	Short: "delete statistics for login",
	Run: func(cmd *cobra.Command, args []string) {
		if login == "" {
			log.Fatal(errors.New("login is not specified"))
			return
		}

		r, err := client.DeleteLoginStats(ctx, &proto.DeleteLoginStatsRequest{Login: login})
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Is success -", r.IsSuccess)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteLoginStatsCmd)

	deleteLoginStatsCmd.PersistentFlags().StringVar(
		&login,
		"login",
		"",
		"Login, for example \"192.168.0.1\"",
	)
}
