package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	StartScheduler()

	router.POST("/upload", func(c *gin.Context) {
		employee := c.PostForm("employee")
		file, err := c.FormFile("summary")
		if err != nil {
			c.String(http.StatusBadRequest, "No file uploaded")
			return
		}
		filePath := "uploads/" + employee + "_" + file.Filename
		c.SaveUploadedFile(file, filePath)
		c.String(http.StatusOK, "Thanks %s! File saved as %s", employee, filePath)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Employee Summary Server Running!")
	})

	router.Run(":8080")
}
