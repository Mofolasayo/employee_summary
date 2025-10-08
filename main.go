package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "employee-summary/docs"
)

var (
	summaries = make(map[string]string)
	grades    = make(map[string]string)
	mu        sync.Mutex
)

func uploadHandler(c *gin.Context) {
	employee := c.PostForm("employee")
	file, err := c.FormFile("summary")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	os.MkdirAll("uploads", os.ModePerm)
	filePath := fmt.Sprintf("uploads/%s_%s", employee, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	SaveSubmission(employee, filePath)

	go func() {
		files, _ := os.ReadDir("uploads")
		combined := ""
		for _, f := range files {
			content, _ := os.ReadFile("uploads/" + f.Name())
			combined += string(content) + "\n\n"
		}

		if combined == "" {
			fmt.Println("No reports found to summarize.")
			return
		}

		summary := SummarizeText(combined)
		evaluation := EvaluateEmployees()

		finalReport := fmt.Sprintf("%s\n\n%s", summary, evaluation)
		SendEmail(finalReport)

		fmt.Println("Email sent successfully with report:\n", finalReport)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File from %s uploaded successfully. Summary will be generated and emailed.", employee),
	})
}

func getAllSummaries(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	response := gin.H{}
	for emp, sum := range summaries {
		response[emp] = gin.H{
			"summary": sum,
			"grade":   grades[emp],
		}
	}
	c.JSON(http.StatusOK, response)
}

func main() {
	_ = godotenv.Load()
	router := gin.Default()

	StartScheduler()

	router.POST("/upload", uploadHandler)
	router.GET("/summaries", getAllSummaries)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Server running on http://localhost:8080 (or your ngrok https URL)")
	router.Run(":8080")
}

func GradePerformance(summary string) string {
	if strings.Contains(strings.ToLower(summary), "excellent") {
		return "A"
	}
	if strings.Contains(strings.ToLower(summary), "improve") {
		return "C"
	}
	return "B"
}
