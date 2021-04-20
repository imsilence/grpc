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
	"google.golang.org/grpc/metadata"
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

		md := metadata.Pairs("type", "multipart")
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		stream, err := client.Multipart(ctx, opts...)
		if err != nil {
			return err
		}

		log.Print(stream.Header())

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

		log.Print(stream.Trailer())
		return nil
	},
}

func init() {
	clientCmd.AddCommand(multipartClientCmd)
}
