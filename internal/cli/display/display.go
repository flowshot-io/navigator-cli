package display

import (
	"strings"

	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/olekukonko/tablewriter"
)

func Assets(assets ...*commandservice.Asset) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Query Results.")
	table.SetHeader([]string{"ID", "Name", "CreatedAt", "UpdatedAt"})

	for _, result := range assets {
		row := []string{
			result.AssetID,
			result.Name,
			result.CreatedAt,
			result.UpdatedAt,
		}
		table.Append(row)
	}

	table.Render()

	return tableString.String()
}

func Files(assets ...*fileservice.File) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Query Results.")
	table.SetHeader([]string{"ID", "AssetID", "Key", "Status", "UpdatedAt"})

	for _, result := range assets {
		row := []string{
			result.FileID,
			result.AssetID,
			result.Key,
			result.Status,
			result.UploadedAt,
		}
		table.Append(row)
	}

	table.Render()

	return tableString.String()
}
