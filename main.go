package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Start weekly scheduler (still runs in background)
	StartScheduler()

	router.POST("/upload", func(c *gin.Context) {
		employee := c.PostForm("employee")
		file, err := c.FormFile("summary")
		if err != nil {
			c.String(http.StatusBadRequest, "No file uploaded")
			return
		}

		// Ensure uploads folder exists
		os.MkdirAll("uploads", os.ModePerm)

		// Save file
		filePath := "uploads/" + employee + "_" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.String(http.StatusInternalServerError, "Error saving file")
			return
		}

		// Respond immediately
		c.String(http.StatusOK, "Thanks %s! File saved as %s. Generating summary now...", employee, filePath)

		// Immediately run AI summarization + email (for testing)
		go func() {
			files, _ := os.ReadDir("uploads")
			combined := ""
			for _, f := range files {
				content, _ := os.ReadFile("uploads/" + f.Name())
				combined += string(content) + "\n\n"
			}
			if combined == "" {
				fmt.Println("No reports found.")
				return
			}

			summary := SummarizeText(combined)
			fmt.Println("Generated summary:\n", summary)
			SendEmail(summary)
			fmt.Println("✅ Email sent!")
		}()
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Employee Summary Server Running!")
	})

	router.Run(":8080")
}

// package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	router := gin.Default()

// 	StartScheduler()

// 	router.POST("/upload", func(c *gin.Context) {
// 		employee := c.PostForm("employee")
// 		file, err := c.FormFile("summary")
// 		if err != nil {
// 			c.String(http.StatusBadRequest, "No file uploaded")
// 			return
// 		}
// 		filePath := "uploads/" + employee + "_" + file.Filename
// 		c.SaveUploadedFile(file, filePath)
// 		c.String(http.StatusOK, "Thanks %s! File saved as %s", employee, filePath)
// 	})

// 	router.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Employee Summary Server Running!")
// 	})

// 	router.Run(":8080")
// }
