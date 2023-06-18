package factory

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/flowshot-io/navigator-client-go/queryservice/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientFactory interface {
	QueryClient() (queryservice.QueryServiceClient, error)
	CommandClient() (commandservice.CommandServiceClient, error)
	FileClient() (fileservice.FileServiceClient, error)
}

type clientFactory struct {
	storageHost  string
	resourceHost string
	searchHost   string
}

func NewClientFactory() ClientFactory {
	return &clientFactory{
		resourceHost: "localhost:50050",
		storageHost:  "localhost:50051",
		searchHost:   "localhost:50052",
	}
}

func (f *clientFactory) QueryClient() (queryservice.QueryServiceClient, error) {
	conn, err := grpc.Dial(f.searchHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create query services gRPC connection: %w", err)
	}

	return queryservice.NewQueryServiceClient(conn), nil
}

func (f *clientFactory) CommandClient() (commandservice.CommandServiceClient, error) {
	conn, err := grpc.Dial(f.resourceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create command services gRPC connection: %w", err)
	}

	return commandservice.NewCommandServiceClient(conn), nil
}

func (f *clientFactory) FileClient() (fileservice.FileServiceClient, error) {
	conn, err := grpc.Dial(f.storageHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create file services gRPC connection: %w", err)
	}

	return fileservice.NewFileServiceClient(conn), nil
}
