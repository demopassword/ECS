package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 로그 파일 생성
	logFile, err := os.Create("/tmp/app.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Gin 라우터 설정
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logFile)) // Gin 로깅을 로그 파일로 리디렉션

	// 이미지를 제공할 엔드포인트
	router.GET("/image", func(c *gin.Context) {
		// 외부 이미지 URL
		imageURL := "https://images.unsplash.com/photo-1533450718592-29d45635f0a9?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8anBnfGVufDB8fDB8fHww&w=1000&q=80"

		// 이미지 다운로드
		resp, err := http.Get(imageURL)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch image")
			return
		}
		defer resp.Body.Close()

		// 이미지를 클라이언트에 반환
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Header("Content-Length", resp.Header.Get("Content-Length"))
		io.Copy(c.Writer, resp.Body)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	port := 8080
	fmt.Printf("Server is running at :%d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
