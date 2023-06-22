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
	SearchClient() (queryservice.QueryServiceClient, error)
	ResourceClient() (commandservice.CommandServiceClient, error)
	StorageClient() (fileservice.FileServiceClient, error)
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

func (f *clientFactory) SearchClient() (queryservice.QueryServiceClient, error) {
	conn, err := grpc.Dial(f.searchHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create search services gRPC connection: %w", err)
	}

	return queryservice.NewQueryServiceClient(conn), nil
}

func (f *clientFactory) ResourceClient() (commandservice.CommandServiceClient, error) {
	conn, err := grpc.Dial(f.resourceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create resource services gRPC connection: %w", err)
	}

	return commandservice.NewCommandServiceClient(conn), nil
}

func (f *clientFactory) StorageClient() (fileservice.FileServiceClient, error) {
	conn, err := grpc.Dial(f.storageHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create storage services gRPC connection: %w", err)
	}

	return fileservice.NewFileServiceClient(conn), nil
}
