package get

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

func (c *Command) NewDownloadCommand() *cobra.Command {
	var outputDir string

	cc := &cobra.Command{
		Use:   "download",
		Short: "Download a file",
		Long:  `Download a file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := c.driver.clientFactory.StorageClient()
			if err != nil {
				return err
			}

			outputPath := filepath.Join(outputDir, c.id)
			startByte := int64(0)

			// Try to resume download if the file exists
			if fileInfo, err := os.Stat(outputPath); err == nil {
				startByte = fileInfo.Size() + 1
			}

			// Open the file for appending, creates it if it doesn't exist
			file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("unable to open file for writing: %w", err)
			}
			defer file.Close()

			readStorageClient, err := client.ReadFile(cmd.Context(), &fileservice.ReadFileRequest{
				FileID:    c.id,
				StartByte: startByte,
			})
			if err != nil {
				return fmt.Errorf("unable to read file: %w", err)
			}

			for {
				resp, err := readStorageClient.Recv()
				if err == io.EOF || resp == nil {
					break
				}
				if err != nil {
					return fmt.Errorf("unable to read file chunk: %w", err)
				}

				// Write the received data to the file
				if _, err := file.Write(resp.Data); err != nil {
					return fmt.Errorf("unable to write file chunk: %w", err)
				}
			}

			cmd.Println("File downloaded successfully.")
			return nil
		},
	}

	cc.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory")

	return cc
}
