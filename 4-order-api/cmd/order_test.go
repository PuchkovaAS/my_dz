package main

import (
	"4-order-api/configs"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	phone      = "+79059324762"
	sessionId  = "VFYai-megsBV0M2BI2pAW9_CukwtHTNw"
	authToken  = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6Iis3OTA1OTMyNDc2MiJ9.mOxEtCjoM4v2-3BRnj3SABzp9lRxwYdEMGhND0s0zZc"
	testDBName = "link_test1"
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

func initDb(t *testing.T) *gorm.DB {
	err := godotenv.Load("test.env")
	if err != nil {
		panic(err)
	}

	DSN := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Создаем саму базу данных
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		t.Fatalf("Failed to connect to DB server: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Подключаемся к новой тестовой БД
	testDSN := fmt.Sprintf("%s dbname=%s", DSN, testDBName)
	gormDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	gormDB.AutoMigrate(
		&product.Product{},
		&user.User{},
		&order.Order{},
		&order.OrderProduct{},
	)
	t.Cleanup(func() {
		// Закрываем соединение
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()

		// Удаляем тестовую БД
		db, err := sql.Open("postgres", DSN)
		if err != nil {
			t.Logf("Warning: failed to connect to drop test database: %v", err)
			return
		}
		defer db.Close()

		// Завершаем все соединения с тестовой БД
		_, err = db.Exec(`
            SELECT pg_terminate_backend(pg_stat_activity.pid)
            FROM pg_stat_activity
            WHERE pg_stat_activity.datname = $1
            AND pid <> pg_backend_pid()`, testDBName)
		if err != nil {
			t.Logf("Warning: failed to terminate connections: %v", err)
		}

		// Удаляем БД
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
		if err != nil {
			t.Logf("Warning: failed to drop test database: %v", err)
		}
	})

	return gormDB
}

func getAllProductIds(db *gorm.DB) []uint {
	var ids []uint
	db.Model(&product.Product{}).Select("id").Find(&ids)
	return ids
}

func TestOrderFlow(t *testing.T) {
	// Инициализация БД
	db := initDb(t)
	defer func() {
		if db != nil {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}
	}()

	// Инициализация тестового сервера ДО выполнения подтестов
	conf := loadConfig()

	conf.Db.Dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		testDBName,
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
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
		// Правильный порядок удаления с учетом foreign key constraints
		// 1. Сначала удаляем связанные записи из order_products
		if orderID != 0 {
			if err := db.Unscoped().Where("order_id = ?", orderID).Delete(&order.OrderProduct{}).Error; err != nil {
				t.Errorf("Failed to delete order items: %v", err)
			}

			// 2. Теперь можно удалить сам заказ
			if err := db.Unscoped().Delete(&order.Order{}, orderID).Error; err != nil {
				t.Errorf("Failed to delete order: %v", err)
			}
		}

		// 3. Удаляем пользователя (после удаления его заказов)
		if err := db.Unscoped().Where("phone = ?", phone).Delete(&user.User{}).Error; err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}

		// 4. Удаляем продукты (после удаления связанных order_products)
		productNames := []string{"orange", "banana"}
		if err := db.Unscoped().Where("name IN ?", productNames).Delete(&product.Product{}).Error; err != nil {
			t.Errorf("Failed to delete test products: %v", err)
		}
	})
}
