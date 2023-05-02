package command

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/navigatorservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewCreateCommand() *cobra.Command {
	var name string

	cc := &cobra.Command{
		Use:   "create",
		Short: "Create assets",
		Long:  `Create assets`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &navigatorservice.CreateAssetRequest{
				Name: name,
			}

			client, err := c.driver.clientFactory.NavigatorClient(cmd)
			if err != nil {
				return fmt.Errorf("unable to create navigator client: %w", err)
			}

			resp, err := client.CreateAsset(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to create asset: %w", err)
			}

			cmd.Println("Created asset: ", resp.Id)

			return nil
		},
	}

	cc.Flags().StringVarP(&name, "name", "n", "", "Name for asset")
	cc.MarkFlagRequired("name")

	return cc
}
