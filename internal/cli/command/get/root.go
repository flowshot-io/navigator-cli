package get

import (
	"github.com/flowshot-io/navigator-cli/internal/cli/factory"
	"github.com/spf13/cobra"
)

type Driver struct {
	clientFactory factory.ClientFactory
}

type Command struct {
	driver *Driver
	id     string
}

func NewDriver(clientFactory factory.ClientFactory) *Driver {
	return &Driver{
		clientFactory: clientFactory,
	}
}

func (d *Driver) NewGetCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "get",
		Short: "Get an object.",
		Long:  `Get an object by ID.`,
	}

	cmd := &Command{
		driver: d,
	}

	cc.AddCommand(
		cmd.NewAssetCommand(),
		cmd.NewFileCommand(),
	)

	cc.PersistentFlags().StringVarP(&cmd.id, "id", "i", "", "The ID of the object to get.")
	cc.MarkFlagRequired("id")

	return cc
}
