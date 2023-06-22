package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-cli/internal/cli/display"
	"github.com/flowshot-io/navigator-client-go/queryservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewSearchCommand() *cobra.Command {
	var searchValue string
	var searchType int
	var displayImage bool
	var pageSize int

	cc := &cobra.Command{
		Use:   "search",
		Short: "Search files.",
		Long:  `Search files.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if searchValue == "" {
				return fmt.Errorf("search input is required")
			}

			if searchType > 3 {
				return fmt.Errorf("invalid search type")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &queryservice.SearchFilesRequest{
				Input: searchValue,
				Type:  queryservice.SearchType(searchType),
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
				PageSize:     pageSize,
			})

			cmd.Println(display.SearchResults(resp.Results...))

			return nil
		},
	}

	cc.Flags().StringVarP(&searchValue, "value", "v", "", "Search value")
	cc.Flags().IntVarP(&searchType, "type", "t", 1, "Search type")
	cc.Flags().BoolVarP(&displayImage, "display-image", "d", false, "Display image")
	cc.Flags().IntVarP(&pageSize, "page-size", "p", 4, "Page size")

	return cc
}
