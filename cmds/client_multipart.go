package cmds

import (
	"context"
	"fmt"
	"grpc/bp/hello"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

var multipartClientCmd = &cobra.Command{
	Use:   "multipart",
	Short: "hello multipart client example",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := dial(addr, clientCertFile)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := hello.NewHelloServiceClient(conn)

		opts := make([]grpc.CallOption, 0)
		if useGzip {
			opts = append(opts, grpc.UseCompressor(gzip.Name))
		}

		stream, err := client.Multipart(context.TODO(), opts...)
		if err != nil {
			return err
		}
		for i := 0; i < rand.Intn(20); i++ {
			err := stream.Send(&hello.HelloRequest{Name: fmt.Sprintf("%d.%s", i, time.Now().Format("2006-01-02 15:04:05"))})
			if err != nil {
				return err
			}
		}
		resp, err := stream.CloseAndRecv()
		if err != nil {
			return err
		}

		log.Print(resp.Reply)
		return nil
	},
}

func init() {
	clientCmd.AddCommand(multipartClientCmd)
}
