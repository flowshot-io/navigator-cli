package create

import (
	"fmt"

	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewFileCommand() *cobra.Command {
	var assetID string
	var key string

	cc := &cobra.Command{
		Use:   "file",
		Short: "create/upload a file",
		Long:  `create/upload a file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &fileservice.CreateFileRequest{
				AssetID: assetID,
				Key:     key,
			}

			client, err := c.driver.clientFactory.FileClient()
			if err != nil {
				return err
			}

			resp, err := client.CreateFile(cmd.Context(), request)
			if err != nil {
				return fmt.Errorf("unable to upload file: %w", err)
			}

			cmd.Println("Uploaded file: ", resp.FileID)

			return nil
		},
	}

	cc.Flags().StringVarP(&assetID, "asset", "a", "", "Asset ID to upload to (required)")
	cc.Flags().StringVarP(&key, "key", "k", "", "Key for file (required)")

	cc.MarkFlagRequired("asset")
	cc.MarkFlagRequired("key")

	return cc
}
