package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {

	start := time.Now()

	c.Next()

	latency := time.Since(start)

	requestID, idErr := c.Get("request_id")

	if !idErr {
		fmt.Println("[ERROR] Cannot get request id")
	}

	requestContent := fmt.Sprintf("[INFO] RequestID: %v - Method: %s - Path: %s - Status: %d - latency: %v \n", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)

	fmt.Println(requestContent)

	os.MkdirAll("logs", 0755)

	file, err := os.OpenFile("logs/request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("[ERROR] %#v", err)
		return
	}

	defer file.Close()

	_, fileErr := file.WriteString(requestContent)

	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	fmt.Println("[DEBUG] Append successfully")

}
