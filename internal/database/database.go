package database

import (
	"fmt"
	"log"

	"garage-barbershop/internal/migrations"
	"garage-barbershop/internal/models"

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
func (d *Database) Migrate(modelList ...interface{}) error {
	if d.DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	err := d.DB.AutoMigrate(modelList...)
	if err != nil {
		return fmt.Errorf("ошибка миграции: %v", err)
	}

	// Создаем начальные роли, если переданы модели ролей
	hasRoleModel := false
	for _, model := range modelList {
		if _, ok := model.(*models.Role); ok {
			hasRoleModel = true
			break
		}
	}

	if hasRoleModel {
		if err := CreateInitialRoles(d.DB); err != nil {
			return fmt.Errorf("ошибка создания начальных ролей: %v", err)
		}
		log.Println("✅ Начальные роли созданы успешно")

		// Мигрируем роли существующих пользователей
		if err := migrations.MigrateExistingUserRoles(d.DB); err != nil {
			return fmt.Errorf("ошибка миграции ролей существующих пользователей: %v", err)
		}
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

// CreateInitialRoles создает начальные роли в системе
func CreateInitialRoles(db *gorm.DB) error {
	// Создаем роли, если они не существуют
	roles := []models.Role{
		{
			Name:        "admin",
			DisplayName: "Администратор",
			Description: "Полный доступ к системе",
			IsActive:    true,
			Permissions: `{"users": ["create", "read", "update", "delete"], "barbers": ["create", "read", "update", "delete"], "appointments": ["create", "read", "update", "delete"]}`,
		},
		{
			Name:        "barber",
			DisplayName: "Барбер",
			Description: "Управление записями и профилем",
			IsActive:    true,
			Permissions: `{"appointments": ["create", "read", "update"], "profile": ["read", "update"]}`,
		},
		{
			Name:        "client",
			DisplayName: "Клиент",
			Description: "Запись на услуги",
			IsActive:    true,
			Permissions: `{"appointments": ["create", "read"], "profile": ["read", "update"]}`,
		},
	}

	for _, role := range roles {
		// Проверяем, существует ли роль
		var existingRole models.Role
		err := db.Where("name = ?", role.Name).First(&existingRole).Error
		if err == gorm.ErrRecordNotFound {
			// Роль не существует, создаем
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			// Другая ошибка
			return err
		}
		// Роль уже существует, пропускаем
	}

	return nil
}
