package navigator

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct {
	Host string
}

func Dial(options Options) (navigatorservice.NavigatorServiceClient, error) {
	if options.Host == "" {
		return nil, fmt.Errorf("host is required")
	}

	conn, err := grpc.Dial(options.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create navigator gRPC connection: %w", err)
	}

	return navigatorservice.NewNavigatorServiceClient(conn), nil
}
