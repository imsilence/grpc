package cmds

import (
	"context"
	"fmt"
	"grpc/bp/hello"
	"grpc/services"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var serverCertFile string
var serverKeyFile string
var token string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server example",
	RunE: func(cmd *cobra.Command, args []string) error {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		opts := make([]grpc.ServerOption, 0)
		if serverCertFile != "" || serverKeyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(serverCertFile, serverKeyFile)
			if err != nil {
				return err
			}
			opts = append(opts, grpc.Creds(creds))
		}

		if token != "" {
			valid := func(ctx context.Context) error {
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					return status.Errorf(codes.InvalidArgument, "missing metadata")
				}
				for _, t := range md["authorization"] {
					if t == fmt.Sprintf("Bearer %s", token) {
						return nil
					}
				}
				return status.Errorf(codes.Unauthenticated, "err token")
			}

			opts = append(opts, grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if err := valid(ctx); err != nil {
					return nil, err
				}
				return handler(ctx, req)
			}))

			opts = append(opts, grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				if err := valid(ss.Context()); err != nil {
					return err
				}
				return handler(srv, ss)
			}))

		}
		server := grpc.NewServer(opts...)
		server.RegisterService(&hello.HelloService_ServiceDesc, &services.HelloService{})
		return server.Serve(listener)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&serverKeyFile, "key", "k", "etc/tls/server/server-key.pem", "key file")
	serverCmd.Flags().StringVarP(&serverCertFile, "cert", "c", "etc/tls/server/server.pem", "cert file")
	serverCmd.Flags().StringVarP(&token, "token", "t", "abc123!@#", "token")
}
