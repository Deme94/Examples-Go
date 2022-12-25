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

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Print the full conf of API-REST",
	Long:  `Print the full conf of API-REST. This file (conf.yaml) must be in the root path of this project.`,
	Example: "\n# Logger\n" +
		`logDir: "./log"` + "\n" +
		`logFileName: "logs"` + "\n" +
		`logFileExt: "txt"` + "\n",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := conf.Setup()
		if err != nil {
			return errors.New("\033[31m" + "CONFIGURATION SERVICE DIDN'T WORK" + "\033[0m")
		}
		fmt.Println("# Logger")
		fmt.Printf("logDir: %s\n", conf.Conf.GetString("logDir"))
		fmt.Printf("logFileName: %s\n", conf.Conf.GetString("logFileName"))
		fmt.Printf("logFileExt: %s\n", conf.Conf.GetString("logFileExt"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(confCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// confCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
