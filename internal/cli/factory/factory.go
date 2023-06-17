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
	fileHost    string
	commandHost string
	queryHost   string
}

func NewClientFactory() ClientFactory {
	return &clientFactory{
		queryHost:   "localhost:50053",
		commandHost: "localhost:50054",
		fileHost:    "localhost:50053",
	}
}

func (f *clientFactory) QueryClient() (queryservice.QueryServiceClient, error) {
	conn, err := grpc.Dial(f.queryHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create query services gRPC connection: %w", err)
	}

	return queryservice.NewQueryServiceClient(conn), nil
}

func (f *clientFactory) CommandClient() (commandservice.CommandServiceClient, error) {
	conn, err := grpc.Dial(f.commandHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create command services gRPC connection: %w", err)
	}

	return commandservice.NewCommandServiceClient(conn), nil
}

func (f *clientFactory) FileClient() (fileservice.FileServiceClient, error) {
	conn, err := grpc.Dial(f.fileHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("unable to create file services gRPC connection: %w", err)
	}

	return fileservice.NewFileServiceClient(conn), nil
}
