package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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
	for _, file := range files {
		if file.Size > 10<<20 {
			c.String(http.StatusBadRequest, "File size exceeds 10 MB limit")
			return
		}

		fileExt := filepath.Ext(file.Filename)
		if fileExt != ".txt" && fileExt != ".pdf" {
			c.String(http.StatusBadRequest, "Only .txt files are allowed")
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

		c.String(http.StatusOK, "File uploaded successfully: %s", file.Filename)
	}

}
