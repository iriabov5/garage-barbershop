package models

import (
	"time"

	"gorm.io/gorm"
)

// Role представляет роль в системе
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Name        string `json:"name" gorm:"uniqueIndex;not null"`        // "admin", "barber", "client"
	DisplayName string `json:"display_name"`                            // "Администратор", "Барбер", "Клиент"
	Description string `json:"description"`                             // Описание роли
	IsActive    bool   `json:"is_active" gorm:"default:true"`           // Активна ли роль
	Permissions string `json:"permissions"`                             // JSON с разрешениями
}

// UserRole представляет связь пользователя с ролью (many-to-many)
type UserRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID     uint      `json:"user_id" gorm:"not null;index"`
	RoleID     uint      `json:"role_id" gorm:"not null;index"`
	AssignedBy uint      `json:"assigned_by"` // Кто назначил роль
	AssignedAt time.Time `json:"assigned_at"` // Когда назначена
	IsActive   int       `json:"is_active" gorm:"default:1"` // Активна ли связь (1 = true, 0 = false)

	// Связи
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// RoleAssignmentRequest представляет запрос на назначение роли
type RoleAssignmentRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	RoleID uint   `json:"role_id" binding:"required"`
	Reason string `json:"reason"` // Причина назначения
}

// RoleRemovalRequest представляет запрос на снятие роли
type RoleRemovalRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	RoleID uint   `json:"role_id" binding:"required"`
	Reason string `json:"reason"` // Причина снятия
}

// UserWithRoles представляет пользователя с его ролями
type UserWithRoles struct {
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}
