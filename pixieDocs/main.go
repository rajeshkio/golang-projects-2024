package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"resty.dev/v3"
)

var (
	version   = "1.0.0"
	buildTime = time.Now().Format(time.RFC3339)
	// These will be overridden during build with -ldflags
)

// Use environment variable for service URLs with fallbacks
var GotenbergURL = getEnv("GOTENBERG_URL", "http://gotenberg-service:3000")
var PixiedocsURL = getEnv("PIXIEDOCS_URL", "http://pixiedocs-service:8080")

// Helper function to get environment variable with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var startTime = time.Now()

// Print a nice startup banner with useful information
func printStartupBanner() {
	banner := `
╔════════════════════════════════════════════════════╗
║                                                    ║
║             PDF Processor Service                  ║
║                                                    ║
╚════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
	log.Printf("Starting PDF Processor version %s", version)
	log.Printf("Build time: %s", buildTime)
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("CPU cores: %d", runtime.NumCPU())

	hostname, err := os.Hostname()
	if err == nil {
		log.Printf("Hostname: %s", hostname)
	}

	log.Printf("Environment: %s", getEnv("ENVIRONMENT", "production"))
	log.Println("----------------------------------------------------")
}

func main() {

	printStartupBanner()

	ginMode := getEnv("GIN_MODE", "release")
	gin.SetMode(ginMode)
	log.Printf("Gin mode: %s", ginMode)

	// Log configured services
	log.Printf("Gotenberg service URL: %s", GotenbergURL)
	log.Printf("Pixiedocs service URL: %s", PixiedocsURL)
	r := gin.Default()

	// Serve frontend
	r.Static("/static", "./static")
	r.LoadHTMLFiles("static/index.html")
	log.Printf("Loaded static assets from ./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/health", func(c *gin.Context) {
		// Check if we can connect to Gotenberg
		client := resty.New().SetTimeout(5 * time.Second)
		gotenbergHealthy := true

		resp, err := client.R().Get(GotenbergURL + "/health")
		if err != nil || resp.StatusCode() != 200 {
			gotenbergHealthy = false
			log.Printf("Warning: Gotenberg service health check failed: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "UP",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   version,
			"buildTime": buildTime,
			"uptime":    time.Since(startTime).String(),
			"checks": []gin.H{
				{
					"name":   "application",
					"status": "UP",
				},
				{
					"name":    "gotenberg",
					"status":  map[bool]string{true: "UP", false: "DOWN"}[gotenbergHealthy],
					"details": map[string]string{"url": GotenbergURL},
				},
			},
		})
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
