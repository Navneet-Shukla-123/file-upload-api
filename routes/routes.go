package routes

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const size int64 = 1000000000000000000

// Home render the HTML file
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)

}

// Upload will upload the pdf or txt file to server
func Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing form: %s", err.Error()))
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No file provided")
		return
	}

	start := time.Now()
	for _, file := range files {
		if file.Size > /*10<<20*/ size {
			c.String(http.StatusBadRequest, "File size exceeds 10 MB limit")
			return
		}

		fileExt := filepath.Ext(file.Filename)
		if fileExt != ".txt" && fileExt != ".pdf" {
			c.String(http.StatusBadRequest, "Only .txt and .pdf files are allowed")
			return
		}

		newFile, err := os.Create(file.Filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create file on server")
			return
		}
		defer newFile.Close()

		uploadedFile, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open uploaded file")
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(newFile, uploadedFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to copy file data")
			return
		}

		end := time.Since(start)

		log.Println("Total time taken is ", end)

		c.String(http.StatusOK, "File uploaded successfully: %s", file.Filename)
	}

}

// Will write the file in chunk
func UploadInChunk(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error parsing form: %s", err.Error()))
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No file provided")
		return
	}

	start := time.Now()
	for _, file := range files {
		if file.Size > /*10<<20*/ size {
			c.String(http.StatusBadRequest, "File size exceeds 10 MB limit")
			return
		}

		fileExt := filepath.Ext(file.Filename)
		if fileExt != ".txt" && fileExt != ".pdf" {
			c.String(http.StatusBadRequest, "Only .txt and .pdf files are allowed")
			return
		}

		// Define chunk size
		chunkSize := 1 * 1024 * 1024 * 100 // 100MB chunks

		log.Println("File name is ", file.Filename)

		newFile, err := os.Create(file.Filename)
		if err != nil {
			log.Println("Error in creating the temporary file ", err)
			c.String(http.StatusInternalServerError, "Failed to create temporary file")
			return
		}

		defer newFile.Close()

		uploadedFile, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to open uploaded file")
			return
		}

		fmt.Println("Uploaded file is ", uploadedFile)

		defer uploadedFile.Close()

		// Chunked reading with error handling
		reader := bufio.NewReader(uploadedFile)
		buffer := make([]byte, chunkSize)
		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					// End of file reached
					break
				}
				c.String(http.StatusInternalServerError, fmt.Sprintf("Error reading file: %s", err.Error()))
				return
			}

			fmt.Println("n from the read is ", n)
			// Write the chunk to the temporary file
			_, err = newFile.Write(buffer[:n])
			if err != nil {
				c.String(http.StatusInternalServerError, "Failed to write chunk to file")
				return
			}
		}

		end := time.Since(start)

		log.Println("Total time taken is ", end)

		c.String(http.StatusOK, "File uploaded successfully: %s", file.Filename)
	}
}
