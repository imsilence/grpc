package cmds

import (
	"context"
	"grpc/bp/hello"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var listClientCmd = &cobra.Command{
	Use:   "list",
	Short: "hello list client example",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := dial(addr, clientCertFile)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := hello.NewHelloServiceClient(conn)
		stream, err := client.List(context.TODO(), &hello.HelloRequest{Name: time.Now().Format("2006-01-02 15:04:05")})
		if err != nil {
			return err
		}
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
			log.Print(resp.Reply)
		}
		return nil
	},
}

func init() {
	clientCmd.AddCommand(listClientCmd)
}
