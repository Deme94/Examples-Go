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

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print the full env of API-REST",
	Long:  `Print the full env of API-REST. This file (.env) must be in the root path of this project.`,
	Example: "\n# Server\n" +
		"VERSION=1.0.0\n" +
		"ENVIRONMENT=development\n" +
		"DOMAIN=mydomain.com\n" +
		"PORT=4000\n" +
		"JWT_SECRET=131j2hk31jh23k12j3h1k3hj\n" +
		"STRIPE_KEY=sk_test_askjdshaskdhakdjh\n" +
		"GOOGLE_LOGIN_CLIENT=213123132-sdads1231231asdgd.apps.googleusercontent.com\n" +
		"\n# Databases\n" +
		"POSTGRES_URI=postgres://<user>:<password>@<host>:<port>/<dbName>?sslmode=disable\n" +
		"MONGO_URI=mongodb://<host>:<port>\n",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := conf.Setup()
		if err != nil {
			return errors.New("\033[31m" + "CONFIGURATION SERVICE DIDN'T WORK" + "\033[0m")
		}
		fmt.Println("# Server")
		fmt.Printf("VERSION=%s\n", conf.Env.GetString("VERSION"))
		fmt.Printf("ENVIRONMENT=%s\n", conf.Env.GetString("ENVIRONMENT"))
		fmt.Printf("DOMAIN=%s\n", conf.Env.GetString("DOMAIN"))
		fmt.Printf("PORT=%d\n", conf.Env.GetInt("PORT"))
		fmt.Printf("JWT_SECRET=%s\n", conf.Env.GetString("JWT_SECRET"))
		fmt.Printf("STRIPE_KEY=%s\n", conf.Env.GetString("STRIPE_KEY"))
		fmt.Printf("GOOGLE_LOGIN_CLIENT=%s\n", conf.Env.GetString("GOOGLE_LOGIN_CLIENT"))
		fmt.Printf("\n# Databases")
		fmt.Printf("POSTGRES_URI=%s\n", conf.Env.GetString("POSTGRES_URI"))
		fmt.Printf("MONGO_URI=%s\n", conf.Env.GetString("MONGO_URI"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
