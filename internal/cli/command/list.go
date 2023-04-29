package command

import (
	"strings"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func (c *Command) NewListCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "list",
		Short: "List assets",
		Long:  `List assets`,
		Run: func(cmd *cobra.Command, args []string) {
			request := &navigatorservice.ListAssetsRequest{}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				cmd.PrintErrln("Unable to create navigator client", err)
				return
			}

			resp, err := client.ListAssets(cmd.Context(), request)
			if err != nil {
				cmd.PrintErrln("Unable to list assets", err)
				return
			}

			cmd.Println(renderAssets(resp.Assets))
		},
	}

	return cc
}

func renderAssets(assets []*navigatorservice.Asset) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Query Results.")
	table.SetHeader([]string{"ID", "Name", "Type"})

	for _, result := range assets {
		row := []string{
			result.Id,
			result.Name,
			result.AssetType,
			// result.CreatedAt.AsTime().Format(time.RFC822),
			// result.UpdatedAt.AsTime().Format(time.RFC822),
		}
		table.Append(row)
	}

	table.Render()

	return tableString.String()
}
