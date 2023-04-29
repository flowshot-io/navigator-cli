package command

import (
	"fmt"
	"strings"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var SearchTypes = []string{"nearText", "nearObject", "nearImage"}

func (c *Command) NewSearchCommand() *cobra.Command {
	var searchType string
	var searchValue string

	cc := &cobra.Command{
		Use:   "search",
		Short: "Search assets",
		Long:  `Search assets`,
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.SearchAssetsRequest{
				SearchType:  searchType,
				SearchValue: searchValue,
			}

			if !contains(SearchTypes, searchType) {
				cmd.PrintErrln("Invalid search type", searchType)
				return
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.SearchAssets(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to search asset", err)
				return
			}

			cmd.Println(renderSearchResults(resp.Results))
		},
	}

	cc.Flags().StringVarP(&searchValue, "value", "v", "", "Value to search for (required)")
	cc.Flags().StringVarP(&searchType, "type", "t", "nearText", "Type of search")

	cc.MarkFlagRequired("value")

	return cc
}

func renderSearchResults(results []*navigatorservice.SearchResult) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Search Results.")
	table.SetHeader([]string{"ID", "Name", "Certainty", "Distance"})

	for _, result := range results {
		row := []string{
			result.Id,
			result.Name,
			fmt.Sprintf("%.2f", result.Certainty),
			fmt.Sprintf("%.2f", result.Distance),
		}
		table.Append(row)
	}

	table.Render()

	return tableString.String()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
