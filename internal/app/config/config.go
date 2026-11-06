package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/haradoxf-cloud/catalog-service/internal/app/config/section"
)

type Config struct {
	// TODO: Добавьте секции Repository, Processor, Monitor
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load() {
	// TODO: Реализуйте загрузку конфигурации
	// 1. Загрузите .env файл (ошибку можно игнорировать)
	// 2. Спарсите переменные окружения в Root с префиксом "APP"
	// 3. При ошибке парсинга — log.Fatal
	_ = godotenv.Load()

	// 2. Парсим переменные окружения с префиксом APP в Root
	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal(err)
	}
}
