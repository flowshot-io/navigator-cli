package list

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewFileCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "file",
		Short: "List files.",
		Long:  `List files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &fileservice.ListFilesRequest{}

			client, err := c.driver.clientFactory.FileClient()
			if err != nil {
				return err
			}

			resp, err := client.ListFiles(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to list files: %w", err)
			}

			cmd.Println(display.Files(resp.Files...))

			return nil
		},
	}

	return cc
}
