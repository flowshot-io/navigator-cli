package create

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/spf13/cobra"
)

const (
	partSize  = 100 * 1024 * 1024 // 100MB
	concLimit = 5                 // number of concurrent upload threads
)

type part struct {
	data []byte
	num  int64
}

func (c *Command) NewFileCommand() *cobra.Command {
	var assetID string
	var filePath string
	var multipart bool

	cc := &cobra.Command{
		Use:   "file",
		Short: "Create/upload a file",
		Long:  `Create/upload a file`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Check if the provided path is valid
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("file does not exist: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx, cancel := context.WithCancel(cmd.Context())
			// defer cancel()

			ctx := context.Background()

			request := &fileservice.CreateFileRequest{
				AssetID: assetID,
			}

			client, err := c.driver.clientFactory.StorageClient()
			if err != nil {
				return err
			}

			cmd.Println("Creating file...")

			resp, err := client.CreateFile(ctx, request)
			if err != nil {
				return fmt.Errorf("unable to create file: %w", err)
			}

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return fmt.Errorf("unable to get file info: %w", err)
			}

			if fileInfo.Size() <= partSize && !multipart {
				data, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("unable to read file: %w", err)
				}

				cmd.Println("Uploading file...")

				_, err = client.WriteFile(ctx, &fileservice.WriteFileRequest{
					FileID: resp.FileID,
					Data:   data,
				})
				if err != nil {
					return fmt.Errorf("unable to write file: %w", err)
				}
			} else {
				// Open the file for reading
				file, err := os.Open(filePath)
				if err != nil {
					return fmt.Errorf("unable to open file: %w", err)
				}
				defer file.Close()

				cmd.Println("Uploading multipart file...")

				partChan := make(chan part)
				errChan := make(chan error, concLimit)
				partNumbers := make([]int64, 0, 1024) // preallocate a large slice to avoid frequent resizing

				var wg sync.WaitGroup
				wg.Add(concLimit)

				for i := 0; i < concLimit; i++ {
					go func() {
						defer wg.Done()

						for p := range partChan {
							cmd.Println("Uploading part", p.num, "of", fileInfo.Size()/partSize+1)
							_, err := client.WriteMultipart(ctx, &fileservice.WriteMultipartRequest{
								FileID:     resp.FileID,
								PartNumber: p.num,
								Data:       p.data,
							})

							if err != nil {
								errChan <- fmt.Errorf("unable to upload part %d: %w", p.num, err)
								//cancel()
								return
							}
						}
					}()
				}

				var partNumber int64 = 1
				for {
					partData := make([]byte, partSize)
					n, err := file.Read(partData)
					if err == io.EOF {
						break
					}
					if err != nil {
						return fmt.Errorf("unable to read file: %w", err)
					}

					select {
					case partChan <- part{data: partData[:n], num: partNumber}:
						partNumbers = append(partNumbers, partNumber)
						partNumber++
					case err := <-errChan:
						close(partChan)
						return err
					case <-ctx.Done():
						close(partChan)
						return fmt.Errorf("operation cancelled")
					}
				}

				close(partChan)
				wg.Wait()

				select {
				case err := <-errChan:
					cmd.PrintErrln("Error uploading part", err)
					return err
				default:
				}

				cmd.Println("Completing multipart upload...")

				// Complete the multipart upload
				_, err = client.CompleteMultipart(ctx, &fileservice.CompleteMultipartRequest{
					FileID:      resp.FileID,
					PartNumbers: partNumbers,
				})
				if err != nil {
					return fmt.Errorf("unable to complete multipart upload: %w", err)
				}
			}

			cmd.Println("File created:", resp.FileID)

			return nil
		},
	}

	cc.Flags().StringVarP(&assetID, "asset", "a", "", "Asset ID to upload to (required)")
	cc.Flags().StringVarP(&filePath, "path", "p", "", "Path to file to upload (required)")
	cc.Flags().BoolVarP(&multipart, "multipart", "m", false, "Force multipart upload")

	cc.MarkFlagRequired("asset")
	cc.MarkFlagRequired("path")

	return cc
}
