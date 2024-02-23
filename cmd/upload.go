package cmd

import (
	"fmt"
	"os"

	"github.com/ephemeralfiles/eph-beta/pkg/config"
	"github.com/ephemeralfiles/eph-beta/pkg/ephcli"
	"github.com/spf13/cobra"
)

// uploadCmd represents the get command
var uploadCmd = &cobra.Command{
	Use:   "up",
	Short: "upload to ephemeralfiles",
	Long:  `upload to ephemeralfiles`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfigApp()
		err := cfg.LoadConfiguration()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		d := ephcli.NewUploader(cfg.Endpoint, cfg.Token)
		err = d.Upload(fileToUpload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
