package list

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewFileCommand() *cobra.Command {
	var searchInput string
	var searchType int32

	cc := &cobra.Command{
		Use:   "file",
		Short: "List files.",
		Long:  `List files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &fileservice.ListFilesRequest{}

			client, err := c.driver.clientFactory.StorageClient()
			if err != nil {
				return err
			}

			resp, err := client.ListFiles(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to list files: %w", err)
			}

			display := display.NewService(&display.Options{
				DisplayImage: c.displayImage,
				PageSize:     c.pageSize,
			})
			cmd.Println(display.Files(resp.Files...))

			return nil
		},
	}

	cc.Flags().StringVarP(&searchInput, "search", "s", "", "Search input")
	cc.Flags().Int32VarP(&searchType, "type", "t", 0, "Search type")

	return cc
}
