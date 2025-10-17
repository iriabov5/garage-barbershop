# 🏗️ SOLID архитектура в Go

## 📚 Что мы изучили в коммите "Implement SOLID architecture"

### 🔧 Что сделали:
- Создали чистую архитектуру согласно SOLID принципам
- Разделили код на слои: config, database, models, repositories, services, handlers
- Реализовали dependency injection через интерфейсы
- Создали тестируемую и расширяемую структуру

### 📚 Как работает:

#### **1. Single Responsibility Principle (SRP)**

**Каждый пакет имеет одну ответственность:**

```
internal/
├── config/          # Только конфигурация
├── database/        # Только подключение к БД
├── models/          # Только модели данных
├── repositories/    # Только доступ к данным
├── services/        # Только бизнес-логика
└── handlers/        # Только HTTP обработка
```

**Java аналог:**
```java
@Configuration
public class DatabaseConfig { ... }

@Service
public class UserService { ... }

@Repository
public interface UserRepository { ... }

@RestController
public class UserController { ... }
```

#### **2. Open/Closed Principle (OCP)**

**Расширяемость без изменения существующего кода:**

**Go интерфейсы:**
```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
}

// Можно добавить новые методы без изменения интерфейса
type AdvancedUserRepository interface {
    UserRepository
    GetByRole(role string) ([]User, error)
    GetActiveUsers() ([]User, error)
}
```

**Java аналог:**
```java
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByTelegramId(Long telegramId);
}

// Расширение через наследование
public interface AdvancedUserRepository extends UserRepository {
    List<User> findByRole(String role);
    List<User> findActiveUsers();
}
```

#### **3. Liskov Substitution Principle (LSP)**

**Интерфейсы можно заменять реализациями:**

```go
// Интерфейс
type UserRepository interface {
    Create(user *User) error
}

// Реализация 1: PostgreSQL
type postgresUserRepository struct {
    db *gorm.DB
}

func (r *postgresUserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

// Реализация 2: In-memory (для тестов)
type memoryUserRepository struct {
    users []User
}

func (r *memoryUserRepository) Create(user *User) error {
    r.users = append(r.users, *user)
    return nil
}

// Обе реализации можно использовать взаимозаменяемо
```

**Java аналог:**
```java
public interface UserRepository {
    User save(User user);
}

@Component
public class JpaUserRepository implements UserRepository {
    // PostgreSQL реализация
}

@Component
public class InMemoryUserRepository implements UserRepository {
    // In-memory реализация
}
```

#### **4. Interface Segregation Principle (ISP)**

**Маленькие, специфичные интерфейсы:**

```go
// Вместо одного большого интерфейса
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByTelegramID(telegramID int64) (*User, error)
    Update(user *User) error
    Delete(id uint) error
    GetBarbers() ([]User, error)
    GetClients() ([]User, error)
}

// Разделяем на специализированные
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}

type UserQueryRepository interface {
    GetByTelegramID(telegramID int64) (*User, error)
    GetBarbers() ([]User, error)
    GetClients() ([]User, error)
}
```

**Java аналог:**
```java
public interface UserRepository {
    User save(User user);
    Optional<User> findById(Long id);
    User update(User user);
    void deleteById(Long id);
}

public interface UserQueryRepository {
    Optional<User> findByTelegramId(Long telegramId);
    List<User> findBarbers();
    List<User> findClients();
}
```

#### **5. Dependency Inversion Principle (DIP)**

**Зависимости через абстракции:**

```go
// Сервис зависит от интерфейса, а не от конкретной реализации
type UserService struct {
    userRepo UserRepository  // Интерфейс, не конкретный тип
}

func NewUserService(userRepo UserRepository) UserService {
    return UserService{userRepo: userRepo}
}

// В main.go инжектим зависимости
func main() {
    // Создаем конкретную реализацию
    userRepo := repositories.NewUserRepository(db)
    
    // Инжектим в сервис
    userService := services.NewUserService(userRepo)
    
    // Инжектим в хендлер
    userHandler := handlers.NewUserHandler(userService)
}
```

**Java аналог:**
```java
@Service
public class UserService {
    private final UserRepository userRepository;
    
    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }
}

// Spring автоматически инжектит зависимости
@Configuration
public class AppConfig {
    @Bean
    public UserRepository userRepository(JpaRepository jpaRepository) {
        return new JpaUserRepository(jpaRepository);
    }
}
```

### 🎯 Почему так:

#### **1. Тестируемость**

**Go:**
```go
func TestUserService(t *testing.T) {
    // Создаем mock репозиторий
    mockRepo := &MockUserRepository{}
    
    // Инжектим в сервис
    service := NewUserService(mockRepo)
    
    // Тестируем
    err := service.CreateUser(&User{Name: "John"})
    assert.NoError(t, err)
}
```

**Java:**
```java
@Test
public void testUserService() {
    // Создаем mock
    UserRepository mockRepo = mock(UserRepository.class);
    
    // Инжектим в сервис
    UserService service = new UserService(mockRepo);
    
    // Тестируем
    User user = new User("John");
    service.createUser(user);
    
    verify(mockRepo).save(user);
}
```

#### **2. Расширяемость**

**Добавление новой функциональности:**

```go
// Добавляем новый сервис без изменения существующих
type NotificationService interface {
    SendWelcomeEmail(user *User) error
    SendAppointmentReminder(appointment *Appointment) error
}

type notificationService struct {
    emailClient EmailClient
}

func (s *notificationService) SendWelcomeEmail(user *User) error {
    // Реализация
}
```

#### **3. Поддерживаемость**

**Четкое разделение ответственности:**
- **Models** - только структуры данных
- **Repositories** - только доступ к данным
- **Services** - только бизнес-логика
- **Handlers** - только HTTP обработка

### ☕ Go vs Java:

| Принцип | Java | Go |
|---------|------|----| 
| **SRP** | `@Service`, `@Repository` | Отдельные пакеты |
| **OCP** | Наследование, интерфейсы | Интерфейсы, композиция |
| **LSP** | Полиморфизм | Интерфейсы |
| **ISP** | Множественные интерфейсы | Маленькие интерфейсы |
| **DIP** | `@Autowired`, конструкторы | Constructor injection |

## 🏗️ Наша архитектура в деталях

### **1. Config Layer**

**Назначение:** Управление конфигурацией

```go
type Config struct {
    Port        string
    Environment string
    DatabaseURL string
    RedisURL    string
    JWTSecret   string
}

func LoadConfig() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        RedisURL:    os.Getenv("REDIS_URL"),
        JWTSecret:   os.Getenv("JWT_SECRET"),
    }
}
```

**Java аналог:**
```java
@Configuration
@ConfigurationProperties(prefix = "app")
public class AppConfig {
    private String port = "8080";
    private String environment = "development";
    private String databaseUrl;
    private String redisUrl;
    private String jwtSecret;
}
```

### **2. Database Layer**

**Назначение:** Подключение к базе данных

```go
type Database struct {
    DB *gorm.DB
}

func NewDatabase(databaseURL string) (*Database, error) {
    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return &Database{DB: db}, nil
}

func (d *Database) Migrate(models ...interface{}) error {
    return d.DB.AutoMigrate(models...)
}
```

**Java аналог:**
```java
@Configuration
public class DatabaseConfig {
    @Bean
    public DataSource dataSource() {
        return DataSourceBuilder.create()
            .url(databaseUrl)
            .build();
    }
    
    @Bean
    public JpaVendorAdapter jpaVendorAdapter() {
        return new HibernateJpaVendorAdapter();
    }
}
```

### **3. Models Layer**

**Назначение:** Структуры данных

```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time     `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
    
    TelegramID int64  `json:"telegram_id" gorm:"uniqueIndex"`
    Username   string `json:"username"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Role       string `json:"role"`
}
```

**Java аналог:**
```java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue
    private Long id;
    
    @CreationTimestamp
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    private LocalDateTime updatedAt;
    
    @Column(unique = true)
    private Long telegramId;
    
    private String username;
    private String firstName;
    private String lastName;
    private String role;
}
```

### **4. Repositories Layer**

**Назначение:** Доступ к данным

```go
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByTelegramID(telegramID int64) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(user *User) error {
    return r.db.Create(user).Error
}
```

**Java аналог:**
```java
@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByTelegramId(Long telegramId);
}

@Repository
public class UserRepositoryImpl implements UserRepository {
    @Autowired
    private JpaRepository jpaRepository;
    
    @Override
    public User save(User user) {
        return jpaRepository.save(user);
    }
}
```

### **5. Services Layer**

**Назначение:** Бизнес-логика

```go
type UserService interface {
    CreateUser(user *User) error
    GetUserByID(id uint) (*User, error)
    RegisterBarber(telegramID int64, username, firstName, lastName string) (*User, error)
}

type userService struct {
    userRepo UserRepository
}

func (s *userService) CreateUser(user *User) error {
    return s.userRepo.Create(user)
}
```

**Java аналог:**
```java
@Service
public class UserService {
    private final UserRepository userRepository;
    
    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }
    
    public User createUser(User user) {
        return userRepository.save(user);
    }
}
```

### **6. Handlers Layer**

**Назначение:** HTTP обработка

```go
type UserHandler struct {
    userService UserService
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.userService.GetUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
```

**Java аналог:**
```java
@RestController
@RequestMapping("/api/users")
public class UserController {
    private final UserService userService;
    
    public UserController(UserService userService) {
        this.userService = userService;
    }
    
    @GetMapping
    public ResponseEntity<List<User>> getUsers() {
        List<User> users = userService.getUsers();
        return ResponseEntity.ok(users);
    }
}
```

## 🎓 Преимущества SOLID архитектуры

### **1. Тестируемость**
- Каждый слой можно тестировать изолированно
- Легко создавать моки для зависимостей
- Unit тесты для каждого компонента

### **2. Расширяемость**
- Новые функции добавляются без изменения существующего кода
- Легко заменить реализацию интерфейса
- Горизонтальное масштабирование

### **3. Поддерживаемость**
- Четкое разделение ответственности
- Легко найти и исправить баги
- Код понятен новым разработчикам

### **4. Переиспользование**
- Компоненты можно использовать в других проектах
- Интерфейсы обеспечивают совместимость
- Модульная архитектура

## 🚀 Следующие шаги

1. **Реализуем API endpoints** для всех моделей
2. **Добавим валидацию** данных
3. **Создадим middleware** для аутентификации
4. **Напишем тесты** для каждого слоя

---

**Следующий урок:** [HTTP серверы](./08-http.md)
