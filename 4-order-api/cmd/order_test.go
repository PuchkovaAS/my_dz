package main

import (
	"4-order-api/configs"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	phone     = "+79059324762"
	sessionId = "VFYai-megsBV0M2BI2pAW9_CukwtHTNw"
	authToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6Iis3OTA1OTMyNDc2MiJ9.mOxEtCjoM4v2-3BRnj3SABzp9lRxwYdEMGhND0s0zZc"
)

func loadConfig() *configs.Config {
	err := godotenv.Load("test.env")
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &configs.Config{
		Db: configs.DbConfig{
			Dsn: fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				os.Getenv("DB_HOST"),
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_NAME"),
				os.Getenv("DB_PORT"),
				os.Getenv("DB_SSLMODE"),
			),
		},
		Logger: configs.LoggerConfig{
			LogFile: os.Getenv("LOG_FILE"),
		},

		Auth: configs.AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}

func initDb() *gorm.DB {
	err := godotenv.Load("test.env")
	if err != nil {
		panic(err)
	}

	DSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Phone:     phone,
		SessionId: sessionId,
		Code:      3452,
	})

	db.Create(&product.Product{
		Name:        "orange",
		Description: "fruit",
		Price:       1.34,
	})

	db.Create(&product.Product{
		Name:        "banana",
		Description: "fruit",
		Price:       4.34,
	})
}

func getAllProductIds(db *gorm.DB) []uint {
	var ids []uint
	db.Model(&product.Product{}).Select("id").Find(&ids)
	return ids
}

func removeData(db *gorm.DB, orderId uint) {
	db.Unscoped().Where("phone = ?", phone).Delete(&user.User{})
	db.Unscoped().Where("name = ?", "orange").Delete(&product.Product{})
	db.Unscoped().Where("name = ?", "banana").Delete(&product.Product{})
	db.Unscoped().Where("id = ?", orderId).Delete(&order.Order{})
}

func TestCreateOrderSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	productIds := getAllProductIds(db)

	conf := loadConfig()
	testServer := httptest.NewServer(App(conf))

	defer testServer.Close()

	for _, productId := range productIds {
		url := testServer.URL + fmt.Sprintf("/product/%d/buy", productId)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", authToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Errorf(
				"Expected status %d, got %d",
				resp.StatusCode,
				http.StatusCreated,
			)
		}
	}

	url := testServer.URL + "/order"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf(
			"Expected status %d, got %d",
			resp.StatusCode,
			http.StatusCreated,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resData order.OrderFormedResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.OrderID == 0 {
		t.Fatal("OrderId is empty")
	}
	removeData(db, resData.OrderID)
}
