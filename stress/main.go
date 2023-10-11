package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// /v1/stress 경로에 대한 POST 핸들러를 등록합니다.
	r.POST("/v1/stress", StressHandler)

	// /healthcheck 경로에 대한 GET 핸들러를 등록하여 기본 상태 확인 엔드포인트를 만듭니다.
	r.GET("/healthcheck", HealthCheckHandler)

	// 서버를 8080 포트에서 실행합니다.
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func StressHandler(c *gin.Context) {
	// 요청 바디에서 CPU 개수를 읽어옵니다.
	cpuCountStr, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read request body: %s", err.Error())})
		return
	}

	cpuCount, err := strconv.Atoi(string(cpuCountStr))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CPU count in request body"})
		return
	}

	// "stress" 명령어를 실행합니다.
	cmd := exec.Command("stress", "-c", strconv.Itoa(cpuCount), "--timeout", "300s")
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to run stress command: %s", err.Error())})
		return
	}

	// 성공적으로 실행되면 200 OK를 반환합니다.
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("stress 명령어가 %d개의 CPU로 실행되었습니다.", cpuCount)})
}

func HealthCheckHandler(c *gin.Context) {
	// 기본적인 상태 확인 응답을 반환합니다.
	c.JSON(http.StatusOK, gin.H{"message": "서비스가 정상 상태입니다."})
}
