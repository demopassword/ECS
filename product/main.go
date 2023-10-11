package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Price    string `json:"price"`
}

func main() {
	// 시크릿 이름을 명령줄 인수로 입력받음
	var secretName string
	flag.StringVar(&secretName, "secretName", "", "Name of the secret in AWS Secret Manager")
	flag.Parse()

	if secretName == "" {
		log.Fatal("Please provide the -secretName flag with the name of the secret in AWS Secret Manager")
	}

	// AWS 세션 생성
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"), // AWS 리전을 설정하세요.
	}))

	// AWS Secret Manager 클라이언트 생성
	secretsManager := secretsmanager.New(awsSession)

	// MySQL 연결 정보 가져오기
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := secretsManager.GetSecretValue(input)
	if err != nil {
		log.Fatal(err)
	}

	// 시크릿 데이터 파싱
	var secretData map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secretData)
	if err != nil {
		log.Fatal(err)
	}

	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		secretData["username"],
		secretData["password"],
		secretData["host"],
		secretData["port"],
		secretData["dbname"])

	// 로그 파일 열기 또는 생성
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 로그 파일에 기록할 로거 생성
	logger := log.New(logFile, "", 0)

	// MySQL 연결
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// HTTP 라우터 설정
	r := gin.Default()

	// 미들웨어를 사용하여 로그를 작성하는 핸들러 함수 추가
	r.Use(func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		protocol := c.Request.Proto
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		logEntry := fmt.Sprintf("%s - (%s) \"%s %s %s %d %.1f \"%s\"\"\n",
			clientIP, endTime.Format("2006-01-02T15:04:05Z"),
			method, path, protocol, statusCode, float64(latency.Microseconds()), userAgent)

		// 로그를 파일과 표준 출력에 기록
		logger.Print(logEntry)
		fmt.Print(logEntry)
	})

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// /v1/product GET 요청을 처리하는 핸들러 (제품 조회)
	r.GET("/v1/product", func(c *gin.Context) {
		// id 쿼리 매개변수를 가져옵니다.
		id := c.DefaultQuery("id", "")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' query parameter"})
			return
		}

		// 여기에 제품 조회 로직을 추가합니다.

		// 데이터베이스에서 데이터를 가져오는 예시:
		var product Product
		err := db.QueryRow("SELECT id, category, price FROM product WHERE id = ?", id).Scan(&product.ID, &product.Category, &product.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	// /v1/product POST 요청을 처리하는 핸들러 (제품 생성 또는 업데이트)
	r.POST("/v1/product", func(c *gin.Context) {
		var requestData Product
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 여기에 제품 생성 또는 업데이트 로직을 추가합니다.

		// 예시: 데이터베이스에 새로운 제품 정보를 삽입하거나 업데이트합니다.
		_, err := db.Exec(
			"INSERT INTO product (id, category, price) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE category = ?, price = ?",
			requestData.ID,
			requestData.Category,
			requestData.Price,
			requestData.Category,
			requestData.Price,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product created or updated successfully"})
	})

	// 서버 시작
	r.Run(":8080")
}
