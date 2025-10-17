# 🧪 Тестирование в Go

## 📚 Что мы изучили в коммите "Add comprehensive testing framework"

### 🔧 Что сделали:
- Создали полную структуру тестирования (Unit, Integration, E2E)
- Реализовали 18 тестовых сценариев
- Настроили фреймворк для тестирования с testify
- Добавили Makefile для удобного запуска тестов
- Создали моки для изоляции тестов

### 📚 Как работает:

#### **1. Пирамида тестирования**

```
    🔺 E2E Tests (10%) - Полные сценарии
   🔺🔺 Integration Tests (20%) - API + БД
  🔺🔺🔺 Unit Tests (70%) - Бизнес-логика
```

#### **2. Unit Tests - Тестируем бизнес-логику**

**Java аналог:**
```java
@ExtendWith(MockitoExtension.class)
class UserServiceTest {
    @Mock
    private UserRepository userRepository;
    
    @InjectMocks
    private UserService userService;
    
    @Test
    void testCreateUser() {
        // Arrange
        User user = new User("John", "Doe");
        when(userRepository.save(user)).thenReturn(user);
        
        // Act
        User result = userService.createUser(user);
        
        // Assert
        assertThat(result).isNotNull();
        verify(userRepository).save(user);
    }
}
```

**Go реализация:**
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    userService := services.NewUserService(mockRepo)
    
    user := &models.User{
        TelegramID: 12345,
        Username:   "testuser",
        FirstName:  "John",
        LastName:   "Doe",
        Role:       "client",
    }
    
    mockRepo.On("Create", user).Return(nil)
    
    // Act
    err := userService.CreateUser(user)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

**Различия:**
- **Java:** `@Mock`, `@InjectMocks` аннотации
- **Go:** Ручное создание моков, `testify/mock`
- **Java:** `when().thenReturn()`, `verify()`
- **Go:** `On().Return()`, `AssertExpectations()`

#### **3. Integration Tests - Тестируем API**

**Java Spring Boot:**
```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.ANY)
class UserControllerIntegrationTest {
    @Autowired
    private TestRestTemplate restTemplate;
    
    @Test
    void testCreateUser() {
        // Arrange
        User user = new User("John", "Doe");
        
        // Act
        ResponseEntity<User> response = restTemplate.postForEntity(
            "/api/users", user, User.class
        );
        
        // Assert
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);
        assertThat(response.getBody().getName()).isEqualTo("John");
    }
}
```

**Go реализация:**
```go
func (suite *APITestSuite) TestCreateUser_Success(t *testing.T) {
    // Arrange
    userData := map[string]interface{}{
        "telegram_id": 12345,
        "username":    "testuser",
        "first_name":  "John",
        "last_name":   "Doe",
        "role":        "client",
    }
    
    jsonData, _ := json.Marshal(userData)
    
    // Act
    resp, err := http.Post(
        suite.server.URL+"/api/users/create",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    // Assert
    suite.NoError(err)
    suite.Equal(http.StatusCreated, resp.StatusCode)
    
    var response models.User
    err = json.NewDecoder(resp.Body).Decode(&response)
    suite.NoError(err)
    suite.Equal("John", response.FirstName)
}
```

**Различия:**
- **Java:** `@SpringBootTest`, `TestRestTemplate`
- **Go:** `httptest.Server`, `http.Post`
- **Java:** Автоматическая настройка тестовой БД
- **Go:** Ручная настройка SQLite in-memory

#### **4. E2E Tests - Тестируем полные сценарии**

**Java Selenium:**
```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class UserJourneyE2ETest {
    @Test
    void testCompleteUserJourney() {
        // 1. Регистрация барбера
        // 2. Регистрация клиента
        // 3. Создание услуги
        // 4. Запись на услугу
        // 5. Проверка результата
    }
}
```

**Go реализация:**
```go
func (suite *UserJourneyTestSuite) TestCompleteUserJourney(t *testing.T) {
    // 1. Регистрация барбера через API
    barberData := map[string]interface{}{
        "telegram_id": 11111,
        "username":    "ivan_barber",
        "first_name":  "Ivan",
        "last_name":   "Barber",
        "role":        "barber",
    }
    
    jsonData, _ := json.Marshal(barberData)
    resp, err := http.Post(
        suite.server.URL+"/api/users/create",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    suite.NoError(err)
    suite.Equal(http.StatusCreated, resp.StatusCode)
    
    // 2. Регистрация клиента
    // 3. Проверка фильтрации по ролям
    // 4. Проверка сохранения в БД
}
```

### 🎯 Почему так:

#### **1. Изоляция тестов**

**Unit Tests:**
- **Моки** изолируют бизнес-логику
- **Быстрые** (без БД, без HTTP)
- **Стабильные** (не зависят от внешних факторов)

**Integration Tests:**
- **Реальная БД** (SQLite in-memory)
- **HTTP запросы** (httptest.Server)
- **Полный стек** (Handler → Service → Repository → DB)

**E2E Tests:**
- **Полные сценарии** пользователя
- **Критические пути** бизнес-процессов
- **Реальные данные** и взаимодействия

#### **2. Покрытие тестами**

**Unit Tests (70%):**
- Бизнес-логика в сервисах
- Алгоритмы и вычисления
- Валидация данных
- Обработка ошибок

**Integration Tests (20%):**
- API endpoints
- Работа с БД
- HTTP запросы/ответы
- Внешние сервисы

**E2E Tests (10%):**
- Полные пользовательские сценарии
- Критические пути
- Интеграция всех компонентов

### ☕ Go vs Java:

| Аспект | Java | Go |
|--------|------|----| 
| **Фреймворк** | JUnit 5, Mockito | testify, httptest |
| **Моки** | `@Mock`, `@InjectMocks` | Ручные моки |
| **Тестовая БД** | `@DataJpaTest`, H2 | SQLite in-memory |
| **HTTP тесты** | `@WebMvcTest`, `TestRestTemplate` | `httptest.Server` |
| **Запуск** | Maven/Gradle | `go test`, Makefile |
| **Покрытие** | JaCoCo | `go test -cover` |

## 🏗️ Наша структура тестов

### **1. Unit Tests - `tests/unit/`**

**Тестовые сценарии:**
- ✅ `TestUserService_CreateUser` - Создание пользователя
- ✅ `TestUserService_CreateUser_Error` - Ошибка при создании
- ✅ `TestUserService_RegisterBarber` - Регистрация барбера
- ✅ `TestUserService_RegisterClient` - Регистрация клиента
- ✅ `TestUserService_GetUserByID` - Получение по ID
- ✅ `TestUserService_GetUserByID_NotFound` - Пользователь не найден

**Особенности:**
- **Моки** для изоляции
- **Быстрые** (без БД)
- **Стабильные** (детерминированные)

### **2. Integration Tests - `tests/integration/`**

**Тестовые сценарии:**
- ✅ `TestGetUsers_Empty` - Пустой список пользователей
- ✅ `TestCreateUser_Success` - Успешное создание через API
- ✅ `TestCreateUser_InvalidData` - Невалидные данные
- ✅ `TestGetUsers_WithData` - Список с данными
- ✅ `TestGetUsers_ByRole` - Фильтрация по ролям
- ✅ `TestAPIStatus` - Статус API
- ✅ `TestUserService_RegisterBarber_Integration` - Интеграционный тест

**Особенности:**
- **Реальная БД** (SQLite in-memory)
- **HTTP запросы** (httptest.Server)
- **Полный стек** (Handler → Service → Repository → DB)

### **3. E2E Tests - `tests/e2e/`**

**Тестовые сценарии:**
- ✅ `TestCompleteUserJourney` - Полный пользовательский сценарий
- ✅ `TestBarberRegistrationFlow` - Сценарий регистрации барбера
- ✅ `TestClientRegistrationFlow` - Сценарий регистрации клиента
- ✅ `TestErrorHandling` - Обработка ошибок
- ✅ `TestAPIStatus` - Статус API

**Особенности:**
- **Полные сценарии** пользователя
- **Критические пути** бизнес-процессов
- **Реальные данные** и взаимодействия

## 🚀 Команды для тестирования

### **Makefile команды:**

```bash
# Юнит тесты (быстрые)
make test-unit

# Интеграционные тесты (API + БД)
make test-integration

# E2E тесты (полные сценарии)
make test-e2e

# Все тесты
make test-all

# Анализ покрытия
make coverage

# Линтинг
make lint

# Форматирование
make fmt
```

### **Прямые команды Go:**

```bash
# Юнит тесты
go test -v ./tests/unit/... -short

# Интеграционные тесты
go test -v ./tests/integration/... -timeout 30s

# E2E тесты
go test -v ./tests/e2e/... -timeout 60s

# Все тесты
go test -v ./tests/... -timeout 120s

# С покрытием
go test -v ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 🎓 Ключевые концепции Go тестирования

### **1. Тестирование с моками**

```go
// Создаем мок
type MockUserRepository struct {
    mock.Mock
}

// Реализуем методы интерфейса
func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Настраиваем поведение
mockRepo.On("Create", user).Return(nil)

// Проверяем вызовы
mockRepo.AssertExpectations(t)
```

### **2. Тестирование с реальной БД**

```go
// Создаем in-memory БД
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

// Выполняем миграции
err = db.AutoMigrate(&models.User{})

// Очищаем после каждого теста
db.Exec("DELETE FROM users")
```

### **3. HTTP тестирование**

```go
// Создаем тестовый сервер
mux := http.NewServeMux()
mux.HandleFunc("/api/users", handler.GetUsers)
server := httptest.NewServer(mux)

// Выполняем HTTP запросы
resp, err := http.Get(server.URL + "/api/users")

// Проверяем ответ
assert.Equal(t, http.StatusOK, resp.StatusCode)
```

### **4. Test Suites**

```go
type APITestSuite struct {
    suite.Suite
    db     *database.Database
    server *httptest.Server
}

func (suite *APITestSuite) SetupSuite() {
    // Настройка перед всеми тестами
}

func (suite *APITestSuite) TearDownSuite() {
    // Очистка после всех тестов
}

func (suite *APITestSuite) SetupTest() {
    // Настройка перед каждым тестом
}
```

## 🎯 Преимущества нашей стратегии тестирования

### **1. Быстрота**
- **Unit тесты** выполняются за миллисекунды
- **Изоляция** через моки
- **Параллельное** выполнение

### **2. Надежность**
- **Стабильные** тесты (не флакают)
- **Детерминированные** результаты
- **Изолированные** сценарии

### **3. Покрытие**
- **70% Unit** - бизнес-логика
- **20% Integration** - API и БД
- **10% E2E** - полные сценарии

### **4. Поддерживаемость**
- **Четкая структура** тестов
- **Переиспользуемые** компоненты
- **Документированные** сценарии

## 🚀 Следующие шаги

1. **Добавить тесты** для остальных сервисов
2. **Настроить CI/CD** с автоматическими тестами
3. **Добавить тесты** для валидации
4. **Создать тесты** для middleware

---

**Следующий урок:** [Деплой и Docker](./10-deployment.md)
