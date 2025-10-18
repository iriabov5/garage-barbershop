package migrations

import (
	"garage-barbershop/internal/models"

	"gorm.io/gorm"
)

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
