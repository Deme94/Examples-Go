/*
Copyright © 2022 Demetrio Navarro Martínez <deme1994@gmail.com>
*/
package cmd

import (
	"API-REST/services/conf"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Print the version number of API-REST",
	Long:         `Print the version number of API-REST`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := conf.Setup()
		if err != nil {
			return errors.New("\033[31m" + "CONFIGURATION SERVICE DIDN'T WORK" + "\033[0m")
		}
		fmt.Println(conf.Env.GetString("VERSION"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
