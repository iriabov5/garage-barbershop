package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database представляет подключение к базе данных
type Database struct {
	DB *gorm.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(databaseURL string) (*Database, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL не установлен")
	}

	// Настройка логирования GORM в зависимости от окружения
	var gormLogger logger.Interface
	if databaseURL != "" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")

	return &Database{DB: db}, nil
}

// Migrate выполняет миграции базы данных
func (d *Database) Migrate(models ...interface{}) error {
	if d.DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	err := d.DB.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("ошибка миграции: %v", err)
	}

	log.Println("✅ Миграция базы данных выполнена успешно")
	return nil
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	if d.DB == nil {
		return nil
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
