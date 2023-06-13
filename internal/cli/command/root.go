package command

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

func NewDriver() *Driver {
	clientFactory := factory.NewClientFactory()

	return &Driver{
		clientFactory: clientFactory,
	}
}

func NewCommand(d *Driver) *cobra.Command {
	c := &cobra.Command{
		Use:           "navigator",
		Short:         "A command-line tool for Navigator.",
		Long:          `A command-line tool for managing stored assets via Navigator.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	c.SetVersionTemplate("{{.Version}}\n")

	cmd := &Command{
		driver: d,
	}

	createCMD := cmd.NewCreateCommand()
	getCMD := cmd.NewGetCommand()
	uploadCMD := cmd.NewUploadCommand()
	downloadCMD := cmd.NewDownloadCommand()
	listCMD := cmd.NewListCommand()
	deleteCMD := cmd.NewDeleteCommand()
	updateCMD := cmd.NewUpdateCommand()
	searchCMD := cmd.NewSearchCommand()
	importCMD := cmd.NewImportCommand()

	c.AddCommand(
		createCMD,
		getCMD,
		uploadCMD,
		downloadCMD,
		listCMD,
		deleteCMD,
		updateCMD,
		searchCMD,
		importCMD,
	)

	return c
}
