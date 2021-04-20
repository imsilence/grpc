package cmds

import (
	"fmt"
	"grpc/resolver/file"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/resolver"
)

var clientCertFile string
var accessToken string
var useGzip bool
var servers []string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client example",
}

func dial(addr, certFile string) (*grpc.ClientConn, error) {
	opts := make([]grpc.DialOption, 0)

	if len(servers) > 0 {
		// r := manual.NewBuilderWithScheme("scheme")
		r := file.NewBuilderWithScheme("scheme")
		addrs := make([]resolver.Address, len(servers), len(servers)+1)
		for i, srv := range servers {
			addrs[i] = resolver.Address{Addr: srv}
		}
		addrs = append(addrs, resolver.Address{Addr: addr})
		r.InitialState(resolver.State{
			Addresses: addrs,
		})

		addr = fmt.Sprintf("%s:///", r.Scheme())
		opts = append(opts, grpc.WithResolvers(r))

		opts = append(opts, grpc.WithDefaultServiceConfig(`
		{
			"LoadBalancingPolicy":"round_robin"
		}
		`))
	}

	if certFile == "" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(certFile, "server.shadow.com")
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	if accessToken != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{AccessToken: accessToken})))
	}
	return grpc.Dial(addr, opts...)

}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().StringSliceVarP(&servers, "servers", "s", []string{}, "servers addr")
	clientCmd.PersistentFlags().StringVarP(&clientCertFile, "cert", "c", "etc/tls/ca/ca.pem", "cert file")
	clientCmd.PersistentFlags().StringVarP(&accessToken, "token", "t", "abc123!@#", "access token")
	clientCmd.PersistentFlags().BoolVarP(&useGzip, "gzip", "g", false, "gzip compress")
}
