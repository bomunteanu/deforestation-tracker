package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"deforestation/models"

	"deforestation/database"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

var apiKey = "egZCFjVhe1y2Mp3AmEonfYsyLqx0mq0c"

func latLonToTile(lat, lon float64, zoom int) (int, int) {
	latRad := lat * math.Pi / 180.0
	n := math.Pow(2.0, float64(zoom))
	x := int((lon + 180.0) / 360.0 * n)
	y := int((1.0 - math.Log(math.Tan(latRad)+1.0/math.Cos(latRad))/math.Pi) / 2.0 * n)
	return x, y
}

func getSatelliteImageTile(z, x, y int) ([]byte, error) {
	url := fmt.Sprintf("https://api.tomtom.com/map/1/tile/sat/main/%d/%d/%d.jpg?key=%s", z, x, y, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve tile (%d/%d/%d). Status code: %d", z, x, y, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func downloadTiles(latMin, lonMin, latMax, lonMax float64, zoom int) ([][][]byte, error) {
	xMin, yMin := latLonToTile(latMax, lonMin, zoom)
	xMax, yMax := latLonToTile(latMin, lonMax, zoom)

	fmt.Printf("Tile range: x_min=%d, y_max=%d, x_max=%d, y_min=%d\n", xMin, yMax, xMax, yMin)

	var tiles [][][]byte
	for y := yMin; y <= yMax; y++ {
		var rowTiles [][]byte
		for x := xMin; x <= xMax; x++ {
			tileImage, err := getSatelliteImageTile(zoom, x, y)
			if err != nil {
				fmt.Printf("Tile (%d, %d, %d) failed to download: %v\n", zoom, x, y, err)
				continue
			}
			rowTiles = append(rowTiles, tileImage)
		}
		tiles = append(tiles, rowTiles)
	}

	if len(tiles) == 0 {
		return nil, fmt.Errorf("no tiles were downloaded")
	}

	return tiles, nil
}

func stitchTiles(tiles [][][]byte, tileSize int) (image.Image, error) {
	if len(tiles) == 0 {
		return nil, fmt.Errorf("no tiles to stitch")
	}

	numRows := len(tiles)
	numCols := len(tiles[0])

	width := numCols * tileSize
	height := numRows * tileSize

	fmt.Printf("Creating stitched image with dimensions: %dx%d\n", width, height)

	stitchedImage := imaging.New(width, height, image.Transparent)

	for rowIdx, rowTiles := range tiles {
		for colIdx, tileImage := range rowTiles {
			img, err := jpeg.Decode(bytes.NewReader(tileImage))
			if err != nil {
				fmt.Printf("Error processing tile %d, %d: %v\n", rowIdx, colIdx, err)
				continue
			}

			xPos := colIdx * tileSize
			yPos := rowIdx * tileSize

			stitchedImage = imaging.Paste(stitchedImage, img, image.Pt(xPos, yPos))
		}
	}

	return stitchedImage, nil
}

func generateStitchedImage(latMin, lonMin, latMax, lonMax float64, zoom int) (*bytes.Buffer, error) {
	tiles, err := downloadTiles(latMin, lonMin, latMax, lonMax, zoom)
	if err != nil {
		return nil, err
	}

	stitchedImage, err := stitchTiles(tiles, 256)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, stitchedImage, imaging.PNG); err != nil {
		return nil, err
	}

	return buf, nil
}

func GetSatelliteImage(areaID uint) error {
	log.Println("Getting Image..")
	zoom := 15
	db := database.GetDB()

	var area models.Area
	if err := db.First(&area, areaID).Error; err != nil {
		log.Println(err)
		return err
	}

	// Generate the stitched image
	buf, err := generateStitchedImage(area.BottomLeftLat, area.BottomLeftLon, area.TopRightLat, area.TopRightLon, zoom)
	if err != nil {
		log.Printf("Error generating stitched image: %v", err)
		return err
	}

	// Define a directory and filename
	imageDir := "/app/images" // This should match the volume mount point
	imageFilename := fmt.Sprintf("area_%d_%s.png", areaID, time.Now().Format("20060102150405"))
	imagePath := fmt.Sprintf("%s/%s", imageDir, imageFilename)

	// Ensure the directory exists
	if err := createDirIfNotExists(imageDir); err != nil {
		log.Printf("Error ensuring directory exists: %v", err)
		return err
	}

	// Save the image to the specified path
	if err := os.WriteFile(imagePath, buf.Bytes(), 0644); err != nil {
		log.Printf("Error saving stitched image: %v", err)
		return err
	}

	fmt.Printf("Stitched image saved as %s\n", imagePath)

	// Send request to the CV microservice to calculate deforestation
	cvMicroserviceURL := "http://computer-vision:5000/calculate-deforestation/" + imageFilename
	resp, err := http.Get(cvMicroserviceURL)
	if err != nil {
		log.Printf("Error requesting CV microservice: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("CV microservice returned non-200 status: %d", resp.StatusCode)
		return err
	}

	// Parse the response from the CV microservice
	var result struct {
		ForestedArea    float64 `json:"forest_coverage"`
		MaskedImagePath string  `json:"masked_image_path"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading CV microservice response: %v", err)
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error unmarshalling CV microservice response: %v", err)
		return err
	}

	// Update the Area model with the deforested area
	area.DeforestedArea = 100 - result.ForestedArea
	log.Println(area.DeforestedArea)
	if err := db.Save(&area).Error; err != nil {
		log.Printf("Error updating Area model: %v", err)
		return err
	}

	// Create a history record
	history := models.History{
		ImagePath:       imagePath,
		MaskedImagePath: result.MaskedImagePath,
		DeforestedArea:  area.DeforestedArea,
		AreaID:          areaID,
		Date:            time.Now(),
	}

	if err := db.Create(&history).Error; err != nil {
		log.Printf("Error saving history record: %v", err)
		return err
	}

	return nil
}

// CreateDirIfNotExists ensures that a directory exists, creating it if necessary
func createDirIfNotExists(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	return nil
}
