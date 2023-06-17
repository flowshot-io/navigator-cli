package get

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewFileCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "file",
		Short: "Get file",
		Long:  `Get file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &fileservice.GetFileRequest{
				FileID: c.id,
			}

			client, err := c.driver.clientFactory.FileClient()
			if err != nil {
				return err
			}

			resp, err := client.GetFile(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to get asset: %w", err)
			}

			cmd.Println(display.Files(resp.File))

			return nil
		},
	}

	return cc
}
