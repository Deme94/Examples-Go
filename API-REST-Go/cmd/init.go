/*
Copyright © 2022 Demetrio Navarro Martínez <deme1994@gmail.com>
*/
package cmd

import (
	"API-REST/services/conf"
	"API-REST/services/database"
	"log"

	"github.com/spf13/cobra"
)

var resetDB bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database and anything that could be initialized for the first time",
	Long:  `Initialize the database (create db and tables) and anything that could be initialized for the first time`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.LstdFlags | log.Lshortfile) // Set default log flags (print file and line)

		// Conf
		log.Println("Loading configuration service...")
		err := conf.Setup()
		if err != nil {
			log.Fatal("\033[31m"+"CONFIGURATION SERVICE FAILED"+"\033[0m"+" -> ", err)
		}
		log.Println("\033[32m" + "CONFIGURATION SERVICE IS RUNNING" + "\033[0m")

		// DB Init
		if resetDB {
			err = database.Init()
			if err != nil {
				log.Fatal("\033[31m"+"DATABASE CREATION FAILED"+"\033[0m"+" -> ", err)
			}
			log.Println("\033[32m" + "DATABASE WAS CREATED SUCCESSFULLY" + "\033[0m")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolVar(&resetDB, "resetdb", false, "Destroy all tables in DB if exist and create them with no data (or initial data)")
}
