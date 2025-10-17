# 🎯 Основы Go vs Java

## 📚 Что мы изучили в последнем коммите

### 🔧 Что сделали:
- Исправили конфликт пакетов в Go
- Переместили `models.go` в `internal/models/`
- Обновили импорты в `main.go`
- Исправили ссылки на модели

### 📚 Как работает:

#### **1. Пакеты в Go**

**Java:**
```java
package com.company.project.models;
public class User { ... }
```

**Go:**
```go
package models
type User struct { ... }
```

**Ключевые различия:**
- **Java:** `package` = пространство имен
- **Go:** `package` = группа файлов в одной директории

#### **2. Правило Go:**
> **Все файлы в одной директории должны принадлежать одному пакету**

**❌ Неправильно:**
```
project/
├── main.go          # package main
└── models.go        # package models  ← КОНФЛИКТ!
```

**✅ Правильно:**
```
project/
├── main.go                    # package main
└── internal/
    └── models/
        └── models.go          # package models
```

### 🎯 Почему так:

#### **1. Простота компиляции**
- Go компилятор работает с пакетами, а не с отдельными файлами
- Один пакет = один результат компиляции

#### **2. Организация кода**
- Четкое разделение ответственности
- Легко найти связанный код

#### **3. Импорты**
```go
import "garage-barbershop/internal/models"
// Теперь можем использовать: models.User
```

### ☕ Go vs Java:

| Аспект | Java | Go |
|--------|------|----| 
| **Пакеты** | `package com.company.project` | `package models` |
| **Импорты** | `import com.company.project.models.User` | `import "garage-barbershop/internal/models"` |
| **Использование** | `User user = new User()` | `user := &models.User{}` |
| **Компиляция** | JVM байт-код | Нативный код |

## 🏗️ Архитектура проекта

### **Структура (что у нас есть):**

```
garage-barbershop/
├── main.go                     # Точка входа
├── internal/                   # Внутренние пакеты
│   ├── models/                 # Модели данных
│   │   └── models.go           # package models
│   ├── repositories/           # Доступ к данным
│   ├── services/               # Бизнес-логика
│   ├── handlers/               # HTTP обработчики
│   ├── config/                 # Конфигурация
│   └── database/               # Подключение к БД
├── docs/learning/              # Документация
└── go.mod                      # Зависимости
```

### **Аналогия с Java:**

```
src/main/java/com/company/project/
├── Main.java                   # Точка входа
├── models/                     # Модели данных
│   └── User.java              # package com.company.project.models
├── repositories/               # Доступ к данным
├── services/                   # Бизнес-логика
├── controllers/                # HTTP контроллеры
├── config/                     # Конфигурация
└── database/                   # Подключение к БД
```

## 🎓 Ключевые концепции Go

### **1. Структуры vs Классы**

**Java класс:**
```java
public class User {
    private Long id;
    private String name;
    
    public User() {}
    public User(Long id, String name) {
        this.id = id;
        this.name = name;
    }
    
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
}
```

**Go структура:**
```go
type User struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}
```

**Различия:**
- **Java:** Классы с инкапсуляцией, геттеры/сеттеры
- **Go:** Структуры с публичными полями, теги для метаданных

### **2. Методы**

**Java:**
```java
public class UserService {
    public User createUser(User user) {
        return userRepository.save(user);
    }
}
```

**Go:**
```go
type UserService struct {
    userRepo UserRepository
}

func (s *UserService) CreateUser(user *User) error {
    return s.userRepo.Create(user)
}
```

**Различия:**
- **Java:** Методы внутри класса
- **Go:** Методы с получателем (receiver) вне структуры

### **3. Интерфейсы**

**Java:**
```java
public interface UserRepository {
    User save(User user);
    Optional<User> findById(Long id);
}
```

**Go:**
```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
}
```

**Различия:**
- **Java:** Явная реализация интерфейса
- **Go:** Неявная реализация (duck typing)

## 🚀 Следующие шаги

1. **Изучим структуры** подробнее
2. **Разберем методы** и получатели
3. **Поняем интерфейсы** в Go
4. **Реализуем API endpoints**

---

**Следующий урок:** [Структуры vs Классы](./03-structs.md)
