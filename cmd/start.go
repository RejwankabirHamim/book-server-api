/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/RejwankabirHamim/api-book-server/data"
	"github.com/RejwankabirHamim/api-book-server/handler"

	"github.com/spf13/cobra"
)

var (
	port     string
	username string
	password string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This will start the server",
	Long:  `We will specify the port to start the server on and also given credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		data.Users[username] = password
		handler.Caller(port)

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "server port")
	startCmd.Flags().StringVarP(&username, "username", "u", "hamim", "username")
	startCmd.Flags().StringVarP(&password, "password", "s", "1234", "password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
