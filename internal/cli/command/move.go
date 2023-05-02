package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewMoveCommand() *cobra.Command {
	var ID string
	var DestinationPath string

	cc := &cobra.Command{
		Use:   "move",
		Short: "Move asset",
		Long:  `Move asset`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &navigatorservice.MoveAssetRequest{
				Id: ID,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			resp, err := client.MoveAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to list assets: %w", err)
			}

			cmd.Println("Moved asset: ", resp.Id)

			return nil
		},
	}

	cc.Flags().StringVarP(&ID, "id", "i", "", "Asset ID to download")
	cc.Flags().StringVarP(&DestinationPath, "destination-path", "d", "", "Destination path to move asset to")

	cc.MarkFlagRequired("id")
	cc.MarkFlagRequired("destination-path")

	return cc
}
