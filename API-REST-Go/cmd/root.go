/*
Copyright © 2022 Demetrio Navarro Martínez <deme1994@gmail.com>
*/
package cmd

import (
	"API-REST/api-gateway"
	"API-REST/services/conf"
	"API-REST/services/database"
	"API-REST/services/gis"
	"API-REST/services/logger"
	"API-REST/services/mail"
	"API-REST/services/sms"
	"API-REST/services/storage"
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
		err := conf.Setup("", "")
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

		// Mail
		log.Println("Loading mail service...")
		err = mail.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"MAIL SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "MAIL SERVICE IS RUNNING" + "\033[0m")

		// SMS
		log.Println("Loading sms service...")
		err = sms.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"SMS SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "SMS SERVICE IS RUNNING" + "\033[0m")

		// Storage
		log.Println("Loading storage service...")
		err = storage.SetupLocal()
		if err != nil {
			log.Fatal("\033[31m"+"STORAGE SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		err = storage.SetupGCS()
		if err != nil {
			log.Fatal("\033[31m"+"STORAGE SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "STORAGE SERVICE IS RUNNING" + "\033[0m")

		// DB
		log.Println("Loading database service...")
		err = database.SetupPostgres()
		if err != nil {
			log.Fatal("\033[31m"+"DATABASE SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		err = database.SetupMongo()
		if err != nil {
			log.Fatal("\033[31m"+"DATABASE SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "DATABASE SERVICE IS RUNNING" + "\033[0m")

		// GIS MapLibre Martin (Vector tiles server)
		log.Println("Loading GIS service...")
		err = gis.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"MARTIN SERVER IS NOT RUNNING"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "MARTIN SERVER IS RUNNING" + "\033[0m")

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
