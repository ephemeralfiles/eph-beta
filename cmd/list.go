package cmd

import (
	"fmt"
	"os"

	"github.com/ephemeralfiles/eph-beta/pkg/config"
	"github.com/ephemeralfiles/eph-beta/pkg/ephcli"
	"github.com/spf13/cobra"
)

// listCmd represents the get command
var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "list files",
	Long:  `list files`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfigApp()
		err := cfg.LoadConfiguration()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		d := ephcli.NewLister(cfg.Endpoint, cfg.Token)
		err = d.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
