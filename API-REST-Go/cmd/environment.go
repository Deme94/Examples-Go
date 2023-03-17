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

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:          "environment",
	Short:        "Print the environment of API-REST",
	Long:         `Print the environment of API-REST`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := conf.Setup("", "")
		if err != nil {
			return errors.New("\033[31m" + "CONFIGURATION SERVICE DIDN'T WORK" + "\033[0m")
		}
		fmt.Println(conf.Env.GetString("ENVIRONMENT"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(environmentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
