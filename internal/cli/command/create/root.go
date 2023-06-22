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
		Short: "Create an object.",
		Long:  "Create actions for object.",
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
