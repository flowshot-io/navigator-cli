package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/queryservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewSearchCommand() *cobra.Command {
	var searchInput string
	var searchTypeInt int32
	var displayImage bool

	cc := &cobra.Command{
		Use:   "search",
		Short: "Search files.",
		Long:  `Search files.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if searchInput == "" {
				return fmt.Errorf("search input is required")
			}

			if searchTypeInt > 3 {
				return fmt.Errorf("invalid search type")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &queryservice.SearchFilesRequest{
				Input: searchInput,
				Type:  queryservice.SearchType(searchTypeInt),
			}

			client, err := c.driver.clientFactory.SearchClient()
			if err != nil {
				return err
			}

			resp, err := client.SearchFiles(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to search files: %w", err)
			}

			display := display.NewService(&display.Options{
				DisplayImage: displayImage,
			})

			cmd.Println(display.SearchResults(resp.Results...))

			return nil
		},
	}

	cc.Flags().StringVarP(&searchInput, "search", "s", "", "Search input")
	cc.Flags().Int32VarP(&searchTypeInt, "type", "t", 1, "Search type")
	cc.Flags().BoolVarP(&displayImage, "display-image", "d", false, "Display image")

	return cc
}
