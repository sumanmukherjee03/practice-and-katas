package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

func main() {
	describe()
	if len(os.Args) < 2 {
		log.Fatal("Need to pass the directory containing images to the program")
	}

	start := time.Now()
	err := walkFiles(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	fmt.Println("Total time taken : ", elapsed)
}

func walkFiles(rootDirPath string) error {
	err := filepath.Walk(rootDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// checks if it is a file
		if !info.Mode().IsRegular() {
			return nil
		}

		contentType, _ := getFileContentType(path)
		if contentType != "image/jpeg" {
			return nil
		}

		thumbnailImg, err := processImage(path)
		if err != nil {
			return err
		}

		err = saveThumbnail(path, thumbnailImg)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// processImage - takes image file as input
// return pointer to thumbnail image in memory.
func processImage(path string) (*image.NRGBA, error) {
	// load the image from file
	srcImage, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	// scale the image to 100px * 100px
	thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

	return thumbnailImage, nil
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	dir := filepath.Dir(srcImagePath)
	filename := filepath.Base(srcImagePath)
	dstImagePath := filepath.Join(dir, "..", "thumbnail", filename)

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

func describe() {
	str := `
This is an example of an image processing program to generate thumbnails from images.

_____________________
	`
	fmt.Println(str)
}
