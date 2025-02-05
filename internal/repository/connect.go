package repository

import (
	"database/sql"
	"fmt"
	"time"
	"webPractice1/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	_ "github.com/spf13/viper"
)

type Config struct {
	DB         Postgres
	MaxRetries int
	RetryDelay time.Duration
}

func ConfigInicialize() *Config {
	return &Config{
		DB:         Postgres{},
		MaxRetries: 3,
		RetryDelay: 3 * time.Second,
	}
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

func PostgresqlConnect() *sql.DB {
	logger := logger.GetLogger()
	if err := godotenv.Load(".env"); err != nil {
		logger.Error(fmt.Sprintf("godotenv can't load env: %s", err))
	}
	cfg := ConfigInicialize()
	if err := envconfig.Process("db", &cfg.DB); err != nil {
		logger.Error(fmt.Sprintf("envconfig cant parse to struct: %s", err))
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("Open connection to db failed: %s", err))
		connectWithRetry(cfg)
	}

	err = db.Ping()
	if err != nil {
		logger.Error(fmt.Sprintf("ping db failed: %s", err))
		connectWithRetry(cfg)
		db.Close()
	}

	//fmt.Println("Successfully connected!")
	return db
}

func connectWithRetry(cfg *Config) *sql.DB {
	logger := logger.GetLogger()
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		fmt.Printf("Попытка подключения к БД (%d/%d)...\n", i+1, cfg.MaxRetries)

		// Пытаемся подключиться к БД
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Printf("Ошибка при открытии соединения: %v", err)
			logger.Error(fmt.Sprintf("Retry connect to DB faild: %s", err))
		} else if err = db.Ping(); err == nil {
			// Успешное подключение
			fmt.Println("Успешное подключение к базе данных!")
			return db
		}

		// Если подключение не удалось, закрываем его
		if db != nil {
			_ = db.Close()
		}

		fmt.Printf("Не удалось подключиться, ожидаем %v перед повторной попыткой...\n", cfg.RetryDelay)
		time.Sleep(cfg.RetryDelay)
	}

	// Возвращаем ошибку, если не удалось подключиться
	logger.Error(fmt.Sprintf("Retry connect to DB faild: %s", err))
	fmt.Printf("не удалось подключиться к базе данных после %d попыток: %s", cfg.MaxRetries, err)
	return nil
}
