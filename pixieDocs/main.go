package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"resty.dev/v3"
)

var GotenbergURL = "https://pixiedocs.rajesh-kumar.in"

func main() {
	r := gin.Default()

	// Serve frontend
	r.Static("/static", "./static")
	r.LoadHTMLFiles("static/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Convert to PDF
	r.POST("/convert", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
			return
		}
		defer openedFile.Close()

		client := resty.New()
		resp, err := client.R().
			SetDoNotParseResponse(true). // Get raw response
			SetFileReader("files", file.Filename, openedFile).
			Post(GotenbergURL + "/forms/libreoffice/convert")

		if err != nil {
			log.Println("Error making request to Gotenberg:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact conversion service"})
			return
		}

		// Always ensure we close the body when we're done
		defer resp.RawResponse.Body.Close()

		// Read the entire body
		bodyBytes, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Check status code
		if resp.StatusCode() != http.StatusOK {
			// We've already read the body, so we can log it as an error message
			log.Println("Gotenberg conversion failed with status:", resp.Status())
			log.Println("Error body:", string(bodyBytes))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Conversion service returned an error"})
			return
		}

		// Log response size for debugging
		log.Println("PDF response size:", len(bodyBytes), "bytes")

		// Send PDF to client
		c.Data(http.StatusOK, "application/pdf", bodyBytes)
	})

	// Merge PDFs
	r.POST("/merge", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Form parsing failed"})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
			return
		}

		client := resty.New()
		req := client.R().SetDoNotParseResponse(true) // Get raw response

		for _, file := range files {
			openedFile, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
				return
			}
			defer openedFile.Close()

			req.SetFileReader("files", file.Filename, openedFile)
		}

		resp, err := req.Post(GotenbergURL + "/forms/pdf/merge")

		if err != nil {
			log.Println("Error making request to Gotenberg:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact merge service"})
			return
		}

		// Always ensure we close the body when we're done
		defer resp.RawResponse.Body.Close()

		// Read the entire body
		bodyBytes, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Check status code
		if resp.StatusCode() != http.StatusOK {
			// We've already read the body, so we can log it as an error message
			log.Println("Gotenberg merge failed with status:", resp.Status())
			log.Println("Error body:", string(bodyBytes))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Merge service returned an error"})
			return
		}

		// Log response size for debugging
		log.Println("PDF response size:", len(bodyBytes), "bytes")

		// Send PDF to client
		c.Data(http.StatusOK, "application/pdf", bodyBytes)
	})

	// Compress PDF
	r.POST("/compress", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
			return
		}
		defer openedFile.Close()

		client := resty.New()
		resp, err := client.R().
			SetDoNotParseResponse(true). // Get raw response
			SetFileReader("files", file.Filename, openedFile).
			Post(GotenbergURL + "/forms/pdf/compress")

		if err != nil {
			log.Println("Error making request to Gotenberg:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact compression service"})
			return
		}

		// Always ensure we close the body when we're done
		defer resp.RawResponse.Body.Close()

		// Read the entire body
		bodyBytes, err := io.ReadAll(resp.RawResponse.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Check status code
		if resp.StatusCode() != http.StatusOK {
			// We've already read the body, so we can log it as an error message
			log.Println("Gotenberg compression failed with status:", resp.Status())
			log.Println("Error body:", string(bodyBytes))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Compression service returned an error"})
			return
		}

		// Log response size for debugging
		log.Println("PDF response size:", len(bodyBytes), "bytes")

		// Send PDF to client
		c.Data(http.StatusOK, "application/pdf", bodyBytes)
	})

	// Run server
	r.Run(":8080")
}
