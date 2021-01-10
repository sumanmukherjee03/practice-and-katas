package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type processingResult struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

func main() {
	describe()
	if len(os.Args) < 2 {
		log.Fatal("Need to pass the directory containing images to the program")
	}

	start := time.Now()
	err := setupPipeline(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	fmt.Println("Total time taken : ", elapsed)
}

func setupPipeline(rootDirPath string) error {
	done := make(chan struct{}) // create a channel to ensure cancellability of the pipeline
	defer close(done)           // Irrespective of the return path from this func, always make sure all the other goroutines terminate

	// first stage of pipeline
	paths, errs := walkFiles(done, rootDirPath)
	processingResults := processImage(done, paths)

	for res := range processingResults {
		if res.err != nil {
			return res.err // On returning error the done channel is going to get closed and the pipeline will get terminated
		}
		err := saveThumbnail(res.srcImagePath, res.thumbnailImage)
		if err != nil {
			return err // On returning error the done channel is going to get closed and the pipeline will get terminated
		}
	}

	// Remember errs is a buffered channel of 1 returned by the filepathWalk which is why it will not block even if there wasnt a receiver before.
	// Also, reading from that channel just once is enough.
	if err := <-errs; err != nil {
		return err // On returning error the done channel is going to get closed and the pipeline will get terminated
	}

	return nil
}

// It is idiomatic in go to pass the done channel as the first argument
func walkFiles(done <-chan struct{}, rootDirPath string) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error, 1) // make this a buffered channel as you dont want to block the sender goroutine while walking through multiple files

	go func() {
		// Deferred close because irrespective of the return path from this goroutine we want to stop sending paths to process if this goroutine terminates
		// ie, close the out channel even if there is an error or the send is cancelled midway
		defer close(out)

		// Get the error from the file walk and push it to errors channel
		errs <- filepath.Walk(rootDirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err // Return error from this callback and bail. The error will be pushed to the errors channel
			}

			// checks if it is a file
			if !info.Mode().IsRegular() {
				return err // Return error from this callback and bail. The error will be pushed to the errors channel
			}

			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return err // Return error from this callback and bail. The error will be pushed to the errors channel
			}

			// To make the send on the output channel cancellable use the done channel in combination
			// with push to out channel within a select
			select {
			case out <- path:
			case <-done:
				return fmt.Errorf("dir traversal for images was cancelled") // Throw an error if cancelled
			}

			return nil
		})
	}()

	return out, errs
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
func processImage(done <-chan struct{}, paths <-chan string) <-chan *processingResult {
	out := make(chan *processingResult)
	const numThumbnailer = 5
	var wg sync.WaitGroup

	thumbnailer := func(path string) (*image.NRGBA, error) {
		// load the image from file
		srcImage, err := imaging.Open(path)
		if err != nil {
			return nil, err
		}
		thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos) // scale the image to 100px * 100px
		return thumbnailImage, nil
	}

	// Parallelize the process with multiple goroutines doing the image processing
	for i := 0; i < numThumbnailer; i++ {
		wg.Add(1)

		go func() {
			for path := range paths {
				var res processingResult
				thumbnailImage, err := thumbnailer(path)
				if err != nil {
					res = processingResult{
						srcImagePath:   path,
						thumbnailImage: nil,
						err:            err,
					}
					// To make the send on the output channel cancellable use the done channel in combination
					// with push to out channel within a select
					select {
					case out <- &res:
						continue // In case of an error, send the result to the output channel and move on to the next iteration
					case <-done:
						return // bail out if cancelled
					}
				}

				res = processingResult{
					srcImagePath:   path,
					thumbnailImage: thumbnailImage,
					err:            err,
				}

				// To make the send on the output channel cancellable use the done channel in combination
				// with push to out channel within a select
				select {
				case out <- &res:
				case <-done:
					return
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
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
