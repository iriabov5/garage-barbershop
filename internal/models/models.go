package models

import (
	"time"

	"gorm.io/gorm"
)

// User - пользователи системы (барбер и клиенты)
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Основная информация
	TelegramID int64  `json:"telegram_id" gorm:"uniqueIndex"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email" gorm:"uniqueIndex"`

	// Прямая авторизация (без Telegram)
	PasswordHash string `json:"-" gorm:"column:password_hash"` // хеш пароля (не возвращаем в JSON)
	AuthMethod   string `json:"auth_method"`                   // "telegram" или "direct"

	// Роли пользователя (many-to-many через UserRole)
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`

	// Для барбера
	IsActive    bool    `json:"is_active"`   // активен ли барбер
	Specialties string  `json:"specialties"` // специализации (стрижки, бороды, etc)
	Experience  int     `json:"experience"`  // опыт в годах
	Rating      float64 `json:"rating"`      // рейтинг барбера

	// Для клиента
	Preferences string `json:"preferences"` // предпочтения клиента
	Notes       string `json:"notes"`       // заметки о клиенте
}

// Service - услуги барбера
type Service struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Основная информация
	Name        string  `json:"name"`        // название услуги
	Description string  `json:"description"` // описание услуги
	Price       float64 `json:"price"`       // цена
	Duration    int     `json:"duration"`    // длительность в минутах
	IsActive    bool    `json:"is_active"`   // активна ли услуга

	// Связи
	BarberID uint `json:"barber_id" gorm:"not null"`
	Barber   User `json:"barber" gorm:"foreignKey:BarberID"`
}

// Appointment - записи на услуги
type Appointment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Время записи
	DateTime time.Time `json:"datetime" gorm:"not null"`
	Duration int       `json:"duration"` // длительность в минутах

	// Статус записи
	Status string `json:"status"` // "pending", "confirmed", "completed", "cancelled"

	// Связи
	ClientID uint `json:"client_id" gorm:"not null"`
	Client   User `json:"client" gorm:"foreignKey:ClientID"`

	BarberID uint `json:"barber_id" gorm:"not null"`
	Barber   User `json:"barber" gorm:"foreignKey:BarberID"`

	ServiceID uint    `json:"service_id" gorm:"not null"`
	Service   Service `json:"service" gorm:"foreignKey:ServiceID"`

	// Дополнительная информация
	Notes         string  `json:"notes"`          // заметки к записи
	Price         float64 `json:"price"`          // цена услуги
	PaymentStatus string  `json:"payment_status"` // "pending", "paid", "refunded"
}

// WorkingHours - рабочие часы барбера
type WorkingHours struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// День недели (1-7, где 1 = понедельник)
	DayOfWeek int `json:"day_of_week" gorm:"not null"`

	// Время работы
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "18:00"

	// Перерыв
	BreakStart string `json:"break_start"` // "13:00"
	BreakEnd   string `json:"break_end"`   // "14:00"

	// Активен ли этот день
	IsActive bool `json:"is_active"`

	// Связь с барбером
	BarberID uint `json:"barber_id" gorm:"not null"`
	Barber   User `json:"barber" gorm:"foreignKey:BarberID"`
}

// Payment - платежи
type Payment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Сумма и валюта
	Amount   float64 `json:"amount"`   // сумма
	Currency string  `json:"currency"` // валюта (RUB, USD, etc)

	// Статус платежа
	Status string `json:"status"` // "pending", "completed", "failed", "refunded"

	// Способ оплаты
	PaymentMethod string `json:"payment_method"` // "telegram", "card", "cash"

	// Связь с записью
	AppointmentID uint        `json:"appointment_id" gorm:"not null"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`

	// Внешние ID
	ExternalID string `json:"external_id"` // ID в платежной системе
	ReceiptURL string `json:"receipt_url"` // ссылка на чек
}

// Review - отзывы клиентов
type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Оценка и комментарий
	Rating  int    `json:"rating"`  // оценка от 1 до 5
	Comment string `json:"comment"` // комментарий клиента

	// Связи
	ClientID uint `json:"client_id" gorm:"not null"`
	Client   User `json:"client" gorm:"foreignKey:ClientID"`

	BarberID uint `json:"barber_id" gorm:"not null"`
	Barber   User `json:"barber" gorm:"foreignKey:BarberID"`

	AppointmentID uint        `json:"appointment_id" gorm:"not null"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`
}
