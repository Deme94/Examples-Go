/*
Copyright © 2022 Demetrio Navarro Martínez <deme1994@gmail.com>
*/
package cmd

import (
	"API-REST/api-gateway"
	"API-REST/services/conf"
	"API-REST/services/database"
	"API-REST/services/logger"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "API-REST",
	Short: "API-REST template",
	Long:  `This API-REST template supports many databases and some services.`,
	Run: func(cmd *cobra.Command, args []string) {

		log.SetFlags(log.LstdFlags | log.Lshortfile) // Set default log flags (print file and line)

		// Conf
		log.Println("Loading configuration service...")
		err := conf.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"CONFIGURATION SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "CONFIGURATION SERVICE IS RUNNING" + "\033[0m")

		// Logger
		log.Println("Loading logging service...")
		err = logger.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"LOGGING SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "LOGGING SERVICE IS RUNNING" + "\033[0m")

		// DB
		log.Println("Loading database service...")
		err = database.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"DATABASE SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "DATABASE SERVICE IS RUNNING" + "\033[0m")

		// API-Gateway
		log.Println("Loading api-gateway...")
		err = api.Start()
		if err != nil {
			log.Fatal("\033[31m"+"API-GATEWAY FAILED"+"\033[0m"+" -> ", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.API-REST.yaml)")
}
