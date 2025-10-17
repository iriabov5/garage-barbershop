# 🏗️ Структуры vs Классы

## 📚 Что мы изучили в коммите "Implement SOLID architecture"

### 🔧 Что сделали:
- Создали модели данных в `internal/models/models.go`
- Реализовали 6 основных сущностей: User, Service, Appointment, WorkingHours, Payment, Review
- Настроили GORM теги для работы с базой данных
- Применили SOLID принципы к архитектуре

### 📚 Как работает:

#### **1. Структуры в Go**

**Наша модель User:**
```go
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
    Email      string `json:"email"`
    
    // Роль пользователя
    Role string `json:"role"` // "barber" или "client"
}
```

#### **2. Теги (Tags) в Go**

**GORM теги:**
```go
`gorm:"primaryKey"`     // Первичный ключ
`gorm:"uniqueIndex"`    // Уникальный индекс
`gorm:"index"`          // Обычный индекс
`gorm:"not null"`       // NOT NULL
```

**JSON теги:**
```go
`json:"id"`             // Поле в JSON
`json:"telegram_id"`    // Змеиный регистр
`json:"-"`              // Исключить из JSON
```

### 🎯 Почему так:

#### **1. Простота vs Гибкость**

**Java JPA:**
```java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(name = "telegram_id", unique = true)
    private Long telegramId;
    
    @Column(name = "first_name")
    private String firstName;
}
```

**Go GORM:**
```go
type User struct {
    ID         uint   `json:"id" gorm:"primaryKey"`
    TelegramID int64  `json:"telegram_id" gorm:"uniqueIndex"`
    FirstName  string `json:"first_name"`
}
```

**Различия:**
- **Java:** Аннотации на полях, больше кода
- **Go:** Теги в одной строке, компактно

#### **2. Автоматические поля**

**Go GORM автоматически добавляет:**
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time      // Автоматически
    UpdatedAt time.Time      // Автоматически  
    DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
}
```

**Java JPA:**
```java
@Entity
public class User {
    @Id
    @GeneratedValue
    private Long id;
    
    @CreationTimestamp
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    private LocalDateTime updatedAt;
    
    @Column(name = "deleted_at")
    private LocalDateTime deletedAt;
}
```

### ☕ Go vs Java:

| Аспект | Java | Go |
|--------|------|----| 
| **Объявление** | `public class User` | `type User struct` |
| **Поля** | `private Long id` | `ID uint` (публичные) |
| **Аннотации** | `@Entity`, `@Id` | Теги в бэктиках |
| **Конструкторы** | `public User() {}` | Автоматически |
| **Геттеры/Сеттеры** | Обязательны | Не нужны |
| **Наследование** | `extends` | Нет (композиция) |

## 🏗️ Наши модели в деталях

### **1. User (Пользователи)**

**Назначение:** Барберы и клиенты системы

**Java аналог:**
```java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue
    private Long id;
    
    @Column(unique = true)
    private Long telegramId;
    
    private String username;
    private String firstName;
    private String lastName;
    private String role; // "barber" или "client"
    
    // Для барбера
    private Boolean isActive;
    private String specialties;
    private Integer experience;
    private Double rating;
    
    // Для клиента
    private String preferences;
    private String notes;
}
```

**Go реализация:**
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time     `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
    
    // Основная информация
    TelegramID int64  `json:"telegram_id" gorm:"uniqueIndex"`
    Username   string `json:"username"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Phone      string `json:"phone"`
    Email      string `json:"email"`
    Role       string `json:"role"` // "barber" или "client"
    
    // Для барбера
    IsActive    bool    `json:"is_active"`
    Specialties string  `json:"specialties"`
    Experience  int     `json:"experience"`
    Rating      float64 `json:"rating"`
    
    // Для клиента
    Preferences string `json:"preferences"`
    Notes       string `json:"notes"`
}
```

### **2. Service (Услуги)**

**Назначение:** Услуги, которые предоставляет барбер

**Связи:**
- Принадлежит барберу (BarberID)
- Имеет цену и длительность

### **3. Appointment (Записи)**

**Назначение:** Записи клиентов на услуги

**Связи:**
- Клиент (ClientID)
- Барбер (BarberID)  
- Услуга (ServiceID)
- Время и статус

### **4. WorkingHours (Рабочие часы)**

**Назначение:** Расписание работы барбера

**Особенности:**
- День недели (1-7)
- Время начала/окончания
- Перерыв
- Активность дня

### **5. Payment (Платежи)**

**Назначение:** Оплата услуг

**Связи:**
- Связан с записью (AppointmentID)
- Статус платежа
- Способ оплаты

### **6. Review (Отзывы)**

**Назначение:** Отзывы клиентов о работе барбера

**Связи:**
- Клиент (ClientID)
- Барбер (BarberID)
- Запись (AppointmentID)
- Рейтинг и комментарий

## 🎓 Ключевые концепции Go

### **1. Композиция vs Наследование**

**Java (наследование):**
```java
public class BaseEntity {
    private Long id;
    private LocalDateTime createdAt;
}

public class User extends BaseEntity {
    private String name;
}
```

**Go (композиция):**
```go
type BaseEntity struct {
    ID        uint      `gorm:"primaryKey"`
    CreatedAt time.Time
}

type User struct {
    BaseEntity
    Name string
}
```

### **2. Указатели vs Ссылки**

**Java:**
```java
User user = new User();
user.setName("John");
```

**Go:**
```go
user := &User{}  // Указатель
user.Name = "John"

// Или
user := User{}   // Значение
user.Name = "John"
```

### **3. Инициализация**

**Java:**
```java
User user = new User();
user.setId(1L);
user.setName("John");
```

**Go:**
```go
// Способ 1: Пустая структура
user := &User{}

// Способ 2: С значениями
user := &User{
    ID:   1,
    Name: "John",
}

// Способ 3: Только некоторые поля
user := &User{
    Name: "John",
    // ID будет 0 (значение по умолчанию)
}
```

## 🚀 Следующие шаги

1. **Изучим методы** и получатели
2. **Разберем интерфейсы** в Go
3. **Реализуем репозитории** для работы с БД
4. **Создадим сервисы** для бизнес-логики

---

**Следующий урок:** [Методы и функции](./04-methods.md)
