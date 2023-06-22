package display

import (
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/flowshot-io/navigator-client-go/commandservice/v1"
	"github.com/flowshot-io/navigator-client-go/fileservice/v1"
	"github.com/flowshot-io/navigator-client-go/queryservice/v1"
	"github.com/nfnt/resize"
	"github.com/olekukonko/tablewriter"
)

const (
	ImageHost       string = "localhost:50054/files"
	SmallImageWidth uint   = 20
	LargeImageWidth uint   = 32
)

type (
	displayFile struct {
		Data  *fileservice.File
		Image *image.Image
		URL   string
	}

	displayResult struct {
		Data  *queryservice.SearchResult
		Image *image.Image
		URL   string
	}

	Options struct {
		ImageHost    string
		DisplayImage bool
		PageSize     int
	}

	Service struct {
		imageHost    string
		displayImage bool
		pageSize     int
	}
)

func NewService(opts *Options) *Service {
	if opts == nil {
		opts = &Options{}
	}

	if opts.ImageHost == "" {
		opts.ImageHost = ImageHost
	}

	return &Service{
		imageHost:    opts.ImageHost,
		displayImage: opts.DisplayImage,
		pageSize:     opts.PageSize,
	}
}

func (s *Service) SearchResults(results ...*queryservice.SearchResult) string {
	if len(results) == 0 {
		return "No results found."
	}

	if s.pageSize == 0 {
		s.pageSize = len(results)
	}

	var displayResults []*displayResult
	for _, result := range results[:s.pageSize] {
		url := s.getImageURL(result.Key)

		var image *image.Image
		if strings.HasPrefix(result.Mime, "image/") && s.displayImage {
			var err error
			image, err = getImageFromURL(url)
			if err != nil {
				fmt.Printf("Failed to get image from url: %s\n", err)
			}
		}

		displayResults = append(displayResults, &displayResult{
			Data:  result,
			Image: image,
			URL:   url,
		})
	}

	return s.RenderSearchResults(displayResults...)
}

func (s *Service) Assets(assets ...*commandservice.Asset) string {
	return Assets(assets...)
}

func (s *Service) Files(files ...*fileservice.File) string {
	if len(files) == 0 {
		return "No files found."
	}

	if s.pageSize == 0 {
		s.pageSize = len(files)
	}

	var displayFiles []*displayFile
	for _, file := range files[:s.pageSize] {
		url := s.getImageURL(file.Key)

		var image *image.Image
		if strings.HasPrefix(file.Mime, "image/") && s.displayImage {
			var err error
			image, err = getImageFromURL(url)
			if err != nil {
				fmt.Printf("Failed to get image from url: %s\n", err)
			}
		}

		displayFiles = append(displayFiles, &displayFile{
			Data:  file,
			Image: image,
			URL:   url,
		})
	}

	return s.RenderFiles(displayFiles...)
}

func Assets(assets ...*commandservice.Asset) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Query Results.")
	table.SetHeader([]string{"ID", "Name", "CreatedAt", "UpdatedAt"})
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	var rows [][]string
	for _, result := range assets {
		rows = append(rows, []string{
			result.AssetID,
			result.Name,
			result.CreatedAt,
			result.UpdatedAt,
		})
	}

	table.AppendBulk(rows)
	table.Render()

	return tableString.String()
}

func (s *Service) RenderSearchResults(results ...*displayResult) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	headers := []string{"ID", "Description", "Certainty", "Distance", "URL"}

	if s.displayImage {
		headers = append(headers, "Image")
	}

	table.SetCaption(true, "Search Results.")
	table.SetHeader(headers)
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	if s.pageSize == 0 {
		s.pageSize = len(results)
	}

	var fileID string
	var description string
	var certainty string
	var distance string
	var rows [][]string

	for _, result := range results {
		if result.Data != nil {
			fileID = result.Data.Id
			description = result.Data.Description
			certainty = fmt.Sprintf("%.2f", result.Data.Certainty)
			distance = fmt.Sprintf("%.2f", result.Data.Distance)
		}

		data := []string{
			fileID,
			description,
			certainty,
			distance,
			result.URL,
		}

		if s.displayImage {
			imgStr := "Image not found"
			if result.Image != nil {
				imgStr = renderImageToBlocks(*result.Image, SmallImageWidth)
			}

			data = append(data, imgStr)
		}

		rows = append(rows, data)
	}

	table.AppendBulk(rows)
	table.Render()

	return tableString.String()
}

func (s *Service) RenderFiles(files ...*displayFile) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	headers := []string{"ID", "AssetID", "Status", "URL"}

	if s.displayImage {
		headers = append(headers, "Image")
	}

	table.SetCaption(true, "Query Results.")
	table.SetHeader(headers)
	table.SetRowLine(true)
	table.SetRowSeparator("-")

	var imageWidth = SmallImageWidth
	var fileID string
	var assetID string
	var status string
	var rows [][]string

	totalFiles := len(files)
	if totalFiles == 1 {
		imageWidth = LargeImageWidth
	}

	for _, result := range files {
		if result.Data != nil {
			fileID = result.Data.FileID
			assetID = result.Data.AssetID
			status = result.Data.Status.String()
		}

		data := []string{
			fileID,
			assetID,
			status,
			result.URL,
		}

		if s.displayImage {
			imgStr := "Image not found"
			if result.Image != nil {
				imgStr = renderImageToBlocks(*result.Image, imageWidth)
			}

			data = append(data, imgStr)
		}

		rows = append(rows, data)
	}

	table.AppendBulk(rows)
	table.Render()

	return tableString.String()
}

func (s *Service) getImageURL(key string) string {
	return fmt.Sprintf("http://%s/%s", s.imageHost, key)
}

func renderImageToBlocks(img image.Image, maxWidth uint) string {
	// Resize image
	img = resize.Resize(maxWidth, maxWidth*uint(img.Bounds().Dy())/uint(img.Bounds().Dx())*2, img, resize.Bilinear)

	// Create a Builder and use it to write our image to a string
	var builder strings.Builder

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			upper := img.At(x, y)
			r1, g1, b1, _ := upper.RGBA()

			var r2, g2, b2 uint32
			if y+1 < bounds.Max.Y {
				lower := img.At(x, y+1)
				r2, g2, b2, _ = lower.RGBA()
			}

			// Add an ANSI escape code for color before each character
			builder.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄", r2>>8, g2>>8, b2>>8, r1>>8, g1>>8, b1>>8))
		}

		builder.WriteString("\033[0m\n")
	}

	return builder.String()
}

func getImageFromURL(url string) (*image.Image, error) {
	// Perform the HTTP request to get the image data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Decode the image
	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return &img, nil
}
