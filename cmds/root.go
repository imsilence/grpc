package cmds

import (
	"log"

	"github.com/spf13/cobra"
)

var addr string

var rootCmd = &cobra.Command{
	Use:    "grpc",
	Short:  "grpc example",
	Hidden: false,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "d", "127.0.0.1:9999", "地址")
}
