package main

import (
	"4-order-api/configs"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"encoding/json"
	"fmt"
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

var (
	testDB         *gorm.DB
	testServer     *httptest.Server
	createdOrderID uint
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

	db.AutoMigrate(
		&product.Product{},
		&user.User{},
		&order.Order{},
		&order.OrderProduct{},
	)
	return db
}

func getAllProductIds(db *gorm.DB) []uint {
	var ids []uint
	db.Model(&product.Product{}).Select("id").Find(&ids)
	return ids
}

func TestOrderFlow(t *testing.T) {
	// Инициализация БД
	db := initDb()
	defer func() {
		if db != nil {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}
	}()

	// Инициализация тестового сервера ДО выполнения подтестов
	conf := loadConfig()
	testServer := httptest.NewServer(App(conf))
	defer testServer.Close() // Гарантированное закрытие после тестов

	// Переменная для хранения ID заказа
	var orderID uint

	t.Run("InitTestData", func(t *testing.T) {
		// Подготовка тестовых данных
		user := &user.User{
			Phone:     phone,
			SessionId: sessionId,
			Code:      3452,
		}
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		products := []*product.Product{
			{Name: "orange", Description: "fruit", Price: 1.34},
			{Name: "banana", Description: "fruit", Price: 4.34},
		}

		for _, p := range products {
			if err := db.Create(p).Error; err != nil {
				t.Fatalf("Failed to create product %s: %v", p.Name, err)
			}
		}
	})

	t.Run("CreateOrder", func(t *testing.T) {
		// Получаем ID созданных продуктов
		var productIDs []uint
		if err := db.Model(&product.Product{}).Select("id").Find(&productIDs).Error; err != nil {
			t.Fatalf("Failed to get product IDs: %v", err)
		}

		// Тестируем покупку продуктов
		for _, productID := range productIDs {
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("%s/product/%d/buy", testServer.URL, productID),
				nil,
			)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Authorization", authToken)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf(
					"Expected status %d, got %d",
					http.StatusCreated,
					resp.StatusCode,
				)
			}
		}

		// Тестируем оформление заказа
		req, err := http.NewRequest("POST", testServer.URL+"/order", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Authorization", authToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf(
				"Expected status %d, got %d",
				http.StatusCreated,
				resp.StatusCode,
			)
		}

		// Парсим ответ
		var response order.OrderFormedResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.OrderID == 0 {
			t.Fatal("OrderID is empty")
		}
		orderID = response.OrderID
	})

	t.Run("CleanupTestData", func(t *testing.T) {
		// Удаление тестовых данных
		if orderID != 0 {
			if err := db.Unscoped().Delete(&order.Order{}, orderID).Error; err != nil {
				t.Errorf("Failed to delete order: %v", err)
			}
		}

		if err := db.Unscoped().Where("phone = ?", phone).Delete(&user.User{}).Error; err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}

		productNames := []string{"orange", "banana"}
		if err := db.Unscoped().Where("name IN ?", productNames).Delete(&product.Product{}).Error; err != nil {
			t.Errorf("Failed to delete test products: %v", err)
		}
	})
}
