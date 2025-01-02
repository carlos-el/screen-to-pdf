package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"sort"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/jung-kurt/gofpdf"
)

// Manually defined variables:
// Rect points can be easily defined using the front page of the book
var CAPTURE_RECT_UPPER_LEFT = robotgo.Point{X: 278, Y: 204}
var CAPTURE_RECT_LOWER_RIGHT = robotgo.Point{X: 1012, Y: 1351}
var NEXT_BUTTON_LOCATION = robotgo.Point{X: 1168, Y: 791}

const LEFT_MARGIN_MM = 25
const RIGHT_MARGIN_MM = 25
const TOP_MARGIN_MM = 25
const BOTTOM_MARGIN_MM = 25

const FINAL_FILE_PATH = "./output/final.pdf"
const TEMPORAL_SCREENSHOTS_FOLDER_PATH = "./screenshots"
const TIME_TO_START = 5 * time.Second
const SCREENSHOTS_DELAY = 1200 * time.Millisecond

// Actual code
const PDF_WIDTH = 210
const PDF_HEIGHT = 297
const PDF_CONTENT_WIDTH = PDF_WIDTH - LEFT_MARGIN_MM - RIGHT_MARGIN_MM
const PDF_CONTENT_HEIGHT = PDF_HEIGHT - TOP_MARGIN_MM - BOTTOM_MARGIN_MM

var CAPTURE_RECT = robotgo.Rect{
	Point: CAPTURE_RECT_UPPER_LEFT,
	Size: robotgo.Size{
		W: CAPTURE_RECT_LOWER_RIGHT.X - CAPTURE_RECT_UPPER_LEFT.X,
		H: CAPTURE_RECT_LOWER_RIGHT.Y - CAPTURE_RECT_UPPER_LEFT.Y,
	},
}

func main() {
	fmt.Printf("Starting process in %d seconds.\n", int(TIME_TO_START.Seconds()))
	time.Sleep(TIME_TO_START)
	fmt.Println("Starting getting screenshots")

	getAllScreenshots()
	fmt.Println("Finished getting screenshots.")

	aggregateScreenshots()
	fmt.Println("Finished creating PDF file.")
}

func goToNextSlide() {
	// Code for going to next slide/page
	// Can be making click in a position or pressing a keyboard key

	robotgo.MoveClick(NEXT_BUTTON_LOCATION.X, NEXT_BUTTON_LOCATION.Y)
	// err := robotgo.KeyPress("right")
	// if err != nil {
	// 	panic(err)
	// }
}

func encodeToBuffer(img image.Image, buf *bytes.Buffer) error {
	return png.Encode(buf, img) // Convert to PNG for consistent comparison
}

func imagesAreEqual(img1, img2 image.Image) bool {
	if img1 == nil || img2 == nil {
		return false
	}

	// Convert images to byte slices
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	err := encodeToBuffer(img1, buf1)
	if err != nil {
		panic(err)
	}
	err = encodeToBuffer(img2, buf2)
	if err != nil {
		panic(err)
	}

	return bytes.Equal(buf1.Bytes(), buf2.Bytes())
}

func getScreenshot() image.Image {
	img, err := robotgo.CaptureImg(
		CAPTURE_RECT.X, CAPTURE_RECT.Y,
		CAPTURE_RECT.W, CAPTURE_RECT.H,
	)
	if err != nil {
		panic(err)
	}

	return img
}

func saveScreenshot(img image.Image, screenshotNumber int) {
	filePath := fmt.Sprintf("%s/%0*d.png", TEMPORAL_SCREENSHOTS_FOLDER_PATH, 5, screenshotNumber)
	err := robotgo.SavePng(img, filePath)
	if err != nil {
		panic(err)
	}
}

func getAllScreenshots() {
	var previousImg image.Image = nil
	count := 0
	for {
		time.Sleep(SCREENSHOTS_DELAY)
		count++

		img := getScreenshot()
		if imagesAreEqual(img, previousImg) {
			break
		}

		fmt.Printf("\tSaving screenshot number %d.\n", count)
		saveScreenshot(img, count)
		previousImg = img
		goToNextSlide()
	}
}

func getAllScreenshotsFilepaths() []string {
	files, err := os.ReadDir(TEMPORAL_SCREENSHOTS_FOLDER_PATH)
	if err != nil {
		panic(err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	sort.Strings(fileNames)
	return fileNames
}

func aggregateScreenshots() {
	// Create a new PDF instance
	pdf := gofpdf.New("Portrait", "mm", "A4", "")

	// List of image paths to add to the PDF
	imageNames := getAllScreenshotsFilepaths()
	for _, imageName := range imageNames {
		pdf.AddPage()

		// Get image dimensions
		options := gofpdf.ImageOptions{ImageType: "", ReadDpi: true}
		imgPath := fmt.Sprintf("%s/%s", TEMPORAL_SCREENSHOTS_FOLDER_PATH, imageName)

		// Add the image to the PDF
		pdf.ImageOptions(
			imgPath,
			LEFT_MARGIN_MM, TOP_MARGIN_MM,
			PDF_CONTENT_WIDTH, PDF_CONTENT_HEIGHT,
			false, options, 0, "",
		)
	}

	// Save the PDF to a file
	err := pdf.OutputFileAndClose(FINAL_FILE_PATH)
	if err != nil {
		panic(err)
	}
}
