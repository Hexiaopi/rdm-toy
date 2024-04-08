package main

import (
	"log"

	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"

	"github.com/hexiaopi/rdm-toy/internal/config"
	"github.com/hexiaopi/rdm-toy/internal/server"
)

var rootCmd = &cobra.Command{
	Use:   "rdm-toy",
	Short: "rdm-toy is a redis manager tool for toy",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("app run err:%v", err)
	}
}
