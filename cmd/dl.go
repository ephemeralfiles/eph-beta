package cmd

import (
	"fmt"
	"os"

	"github.com/ephemeralfiles/eph-beta/pkg/config"
	"github.com/ephemeralfiles/eph-beta/pkg/ephcli"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var downloadCmd = &cobra.Command{
	Use:   "dl",
	Short: "download from ephemeralfiles",
	Long:  `download from ephemeralfiles`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfigApp()
		err := cfg.LoadConfiguration()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if uuidFileToDownload == "" {
			fmt.Println("uuid is required")
			os.Exit(1)
		}
		d := ephcli.NewDownloader(cfg.Endpoint, cfg.Token)
		err = d.Download(uuidFileToDownload, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
