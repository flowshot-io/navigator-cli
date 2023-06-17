package create

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/factory"
	"github.com/spf13/cobra"
)

type Driver struct {
	clientFactory factory.ClientFactory
}

func NewDriver(clientFactory factory.ClientFactory) *Driver {
	return &Driver{
		clientFactory: clientFactory,
	}
}

type Command struct {
	driver *Driver
}

func (d *Driver) NewCreateCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "create",
		Short: "create an object.",
		Long:  `create an object.`,
	}

	cmd := &Command{
		driver: d,
	}

	cc.AddCommand(
		cmd.NewAssetCommand(),
	)

	return cc
}
