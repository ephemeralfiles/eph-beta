package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configurationFile string
var fileToUpload string
var uuidFileToDownload string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eph",
	Short: "ephemeralfiles command line interface",
	Long:  `ephemeralfiles command line interface`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = os.Getenv("HOME")
		if homedir == "" {
			fmt.Println("warn: $HOME not set")
		}
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// -c option on rootCmd
	rootCmd.PersistentFlags().StringVarP(&configurationFile, "config", "c", filepath.Join(homedir, ".eph.yml"), "configuration file (default is $HOME/.eph.yml))")
	// upload subcommand parameters
	uploadCmd.PersistentFlags().StringVarP(&fileToUpload, "input", "i", "", "file to upload")
	// download subcommand parameters
	downloadCmd.PersistentFlags().StringVarP(&uuidFileToDownload, "input", "i", "", "uuid of file to download")
	// add subcommands
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(listCmd)
}
