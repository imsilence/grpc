package cmds

import (
	"context"
	"fmt"
	"grpc/bp/hello"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var channelClientCmd = &cobra.Command{
	Use:   "channel",
	Short: "hello channel client example",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := dial(addr, clientCertFile)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := hello.NewHelloServiceClient(conn)
		stream, err := client.Channel(context.TODO())
		if err != nil {
			return err
		}
		for i := 0; i < rand.Intn(20); i++ {
			err := stream.Send(&hello.HelloRequest{Name: fmt.Sprintf("%d.%s", i, time.Now().Format("2006-01-02 15:04:05"))})
			if err != nil {
				return err
			}
		}
		err = stream.CloseSend()
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
	clientCmd.AddCommand(channelClientCmd)
}
