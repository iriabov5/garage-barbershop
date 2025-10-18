package migrations

import (
	"log"

	"gorm.io/gorm"
)

// MigrateExistingUserRoles переносит роли из старой системы в новую
func MigrateExistingUserRoles(db *gorm.DB) error {
	log.Println("🔄 Начинаем миграцию ролей существующих пользователей...")

	// 1. Получаем всех пользователей со старой системой ролей
	var users []struct {
		ID   uint   `gorm:"column:id"`
		Role string `gorm:"column:role"`
	}

	// Проверяем, есть ли колонка role в таблице users
	if db.Migrator().HasColumn(&struct{ Role string }{}, "role") {
		err := db.Table("users").Select("id, role").Where("role IS NOT NULL AND role != ''").Find(&users).Error
		if err != nil {
			log.Printf("❌ Ошибка получения пользователей: %v", err)
			return err
		}

		log.Printf("📊 Найдено %d пользователей с ролями для миграции", len(users))

		// 2. Для каждого пользователя назначаем роль в новой системе
		for _, user := range users {
			// Получаем роль по имени
			var role struct {
				ID uint `gorm:"column:id"`
			}
			err := db.Table("roles").Select("id").Where("name = ?", user.Role).First(&role).Error
			if err != nil {
				log.Printf("⚠️ Роль '%s' не найдена для пользователя ID %d: %v", user.Role, user.ID, err)
				continue
			}

			// Проверяем, не назначена ли уже роль
			var count int64
			err = db.Table("user_roles").Where("user_id = ? AND role_id = ? AND is_active = 1", user.ID, role.ID).Count(&count).Error
			if err != nil {
				log.Printf("❌ Ошибка проверки существующей роли: %v", err)
				continue
			}

			if count > 0 {
				log.Printf("✅ Роль '%s' уже назначена пользователю ID %d", user.Role, user.ID)
				continue
			}

			// Назначаем роль
			err = db.Exec(`
				INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at, is_active, created_at, updated_at)
				VALUES (?, ?, ?, NOW(), 1, NOW(), NOW())
			`, user.ID, role.ID, user.ID).Error

			if err != nil {
				log.Printf("❌ Ошибка назначения роли '%s' пользователю ID %d: %v", user.Role, user.ID, err)
				continue
			}

			log.Printf("✅ Роль '%s' успешно назначена пользователю ID %d", user.Role, user.ID)
		}
	} else {
		log.Println("ℹ️ Колонка 'role' не найдена в таблице users - миграция не требуется")
	}

	log.Println("✅ Миграция ролей существующих пользователей завершена")
	return nil
}
