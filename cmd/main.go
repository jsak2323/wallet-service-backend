package main

import (
	"log"

	cron "github.com/btcid/wallet-services-backend-go/cmd/schedulingtask"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd := cobra.Command{
		Use:   "app",
		Short: "Wallet Service API",
		Run: func(*cobra.Command, []string) {
			mainServ()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "run-cron",
		Short: "Run",
		Run: func(cmd *cobra.Command, args []string) {
			cron.Run(args)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "run-cron",
		Short: "Run",
		Run: func(cmd *cobra.Command, args []string) {
			cron.Run(args)
		},
	})

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
