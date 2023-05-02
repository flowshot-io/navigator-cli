package factory

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/navigator"
	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

type ClientFactory interface {
	NavigatorClient(c *cobra.Command) (navigatorservice.NavigatorServiceClient, error)
}

type clientFactory struct {
}

func NewClientFactory() ClientFactory {
	return &clientFactory{}
}

func (f *clientFactory) NavigatorClient(c *cobra.Command) (navigatorservice.NavigatorServiceClient, error) {
	nopts := navigator.Options{Host: "localhost:50052"}

	navigator, err := navigator.Dial(nopts)
	if err != nil {
		return nil, fmt.Errorf("unable to create Navigator client: %w", err)
	}

	return navigator, nil
}
