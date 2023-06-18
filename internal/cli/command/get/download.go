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
			startChunkNumber := int64(1)

			const chunkSize = 5 << 20 // 5 MB

			// Try to resume download if the file exists
			if fileInfo, err := os.Stat(outputPath); err == nil {
				startChunkNumber = fileInfo.Size()/chunkSize + 1
			}

			// Open the file for appending, creates it if it doesn't exist
			file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("unable to open file for writing: %w", err)
			}
			defer file.Close()

			for {
				request := &fileservice.ReadFileRequest{
					FileID:           c.id,
					StartChunkNumber: startChunkNumber,
				}

				readStorageClient, err := client.ReadFile(cmd.Context(), request)
				if err != nil {
					return fmt.Errorf("unable to read file: %w", err)
				}

				resp, err := readStorageClient.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("unable to receive file chunk: %w", err)
				}

				// Write the received data to the file
				if _, err := file.Write(resp.Data); err != nil {
					return fmt.Errorf("unable to write file chunk: %w", err)
				}

				// Increase chunk number for the next request
				startChunkNumber++
			}

			cmd.Println("File downloaded successfully.")
			return nil
		},
	}

	cc.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory")

	return cc
}
