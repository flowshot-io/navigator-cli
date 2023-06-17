package list

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/factory"
	"github.com/spf13/cobra"
)

type Driver struct {
	clientFactory factory.ClientFactory
}

type Command struct {
	driver *Driver
}

func NewDriver(clientFactory factory.ClientFactory) *Driver {
	return &Driver{
		clientFactory: clientFactory,
	}
}

func (d *Driver) NewListCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "list",
		Short: "List objects.",
		Long:  `List objects.`,
	}

	cmd := &Command{
		driver: d,
	}

	cc.AddCommand(
		cmd.NewAssetCommand(),
		cmd.NewFileCommand(),
	)

	return cc
}
