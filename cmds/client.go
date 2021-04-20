package cmds

import (
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var clientCertFile string
var accessToken string
var useGzip bool

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client example",
}

func dial(addr, certFile string) (*grpc.ClientConn, error) {
	opts := make([]grpc.DialOption, 0)
	if certFile == "" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
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
	clientCmd.PersistentFlags().StringVarP(&clientCertFile, "cert", "c", "etc/tls/ca/ca.pem", "cert file")
	clientCmd.PersistentFlags().StringVarP(&accessToken, "token", "t", "abc123!@#", "access token")
	clientCmd.PersistentFlags().BoolVarP(&useGzip, "gzip", "g", false, "gzip compress")
}
