package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := cobra.Command{
		Use:   "app",
		Short: "Wallet Service API",
		Run: func(*cobra.Command, []string) {
			serv()
		},
	}

	// cmd.AddCommand(&cobra.Command{
	// 	Use:   "run-server",
	// 	Short: "Run",
	// 	Run: func(*cobra.Command, []string) {
	// 		serv()
	// 	},
	// })

	// cmd.AddCommand(&cobra.Command{
	// 	Use:   "run-nsq",
	// 	Short: "Run NSQ Consumer",
	// 	Run: func(*cobra.Command, []string) {
	// 		// nsq()
	// 	},
	// })

	// cmd.AddCommand(&cobra.Command{
	// 	Use:   "run-cron",
	// 	Short: "Run Cron",
	// 	Run: func(*cobra.Command, []string) {
	// 		// schedulingTask(1)
	// 	},
	// })

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
