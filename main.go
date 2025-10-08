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
		content, _ := os.ReadFile(filePath)
		summary := SummarizeText(string(content))
		evaluation := EvaluateEmployees()
		finalReport := fmt.Sprintf("Summary for %s:\n%s\n\n%s", employee, summary, evaluation)

		mu.Lock()
		summaries[employee] = summary
		grades[employee] = GradePerformance(summary)
		mu.Unlock()

		SendEmail(finalReport)
		fmt.Println("Email sent successfully with report for", employee)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File from %s uploaded successfully. Summary will be generated and emailed.", employee),
	})
}

func getAllSummaries(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if len(summaries) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No summaries available yet."})
		return
	}

	response := gin.H{}
	for emp, sum := range summaries {
		response[emp] = gin.H{
			"summary": sum,
			"grade":   grades[emp],
		}
	}
	c.JSON(http.StatusOK, response)
}

func getEmployeeSummary(c *gin.Context) {
	employee := c.Param("employee")

	mu.Lock()
	defer mu.Unlock()

	summary, exists := summaries[employee]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No summary found for %s", employee)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employee": employee,
		"summary":  summary,
		"grade":    grades[employee],
	})
}

func main() {
	_ = godotenv.Load()
	router := gin.Default()

	StartScheduler()

	router.POST("/upload", uploadHandler)
	router.GET("/summaries", getAllSummaries)
	router.GET("/summaries/:employee", getEmployeeSummary)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Server running on http://localhost:8080 (or your ngrok https URL)")
	router.Run(":8080")
}

func GradePerformance(summary string) string {
	s := strings.ToLower(summary)
	switch {
	case strings.Contains(s, "excellent"):
		return "A"
	case strings.Contains(s, "good"):
		return "B"
	case strings.Contains(s, "improve"):
		return "C"
	default:
		return "B+"
	}
}
