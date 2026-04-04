# Fitness Journal API

REST API для ведения журнала тренировок с поддержкой персонализированных упражнений, учётом подходов, весов и RPE.

## Стек

- **Go 1.24** (чистый stdlib, минимум зависимостей)
- **PostgreSQL 15** (реляционная БД, транзакции)
- **JWT** (access + refresh tokens, HttpOnly cookies)
- **pgx** (драйвер БД)
- **bcrypt** (хеширование паролей)

## Архитектура

Проект следует чистой архитектуре с разделением по слоям:

```
cmd/app/          — точка входа
internal/
├── app/          — инициализация и запуск сервера
├── config/       — конфигурация (.env)
├── domain/       — доменные модели и ошибки
├── service/      — бизнес-логика
├── storage/      — работа с БД (транзакции, bulk insert)
└── transport/    — HTTP handlers, DTO, middleware
migrations/       — миграции (goose)
```

## Возможности

### Авторизация
- Регистрация / вход / выход
- Access + Refresh токены в HttpOnly cookies
- JWT middleware для защищённых маршрутов

### Упражнения
- CRUD упражнений
- Глобальные и персонализированные упражнения
- Пагинация и фильтрация
- Контроль прав доступа (редактирование только своих)

### Профиль спортсмена
- Получение / обновление / удаление профиля
- Персональные данные (возраст, пол, вес, текущий цикл тренировок)

### Тренировки (WIP)
- Создание тренировки со статусом `in_progress` / `finished`
- Bulk insert подходов (сеты) с проверкой прав на упражнения
- Валидация бизнес-правил (grade, сеты, статус)

## Быстрый старт

### 1. Клонировать и настроить

```bash
git clone <repo-url>
cd fitness
cp .env.example .env
```

Отредактируй `.env`:

```env
DSN=postgres://postgres:postgres@localhost:5432/fitness_journal?sslmode=disable
PORT=8080
JWT_SECRET=your-secret-key-here
REFRESH_JWT_SECRET=your-refresh-secret-here
```

### 2. Запустить PostgreSQL

```bash
docker compose up -d
```

### 3. Применить миграции

```bash
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/fitness_journal?sslmode=disable" up
```

### 4. Запустить сервер

```bash
go run cmd/app/main.go
```

Сервер запустится на `:8080`

## API Endpoints

### Auth

| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| POST | `/auth/register` | Регистрация | ❌ |
| POST | `/auth/login` | Вход | ❌ |
| POST | `/auth/refresh` | Обновление токена | ❌ |
| POST | `/auth/logout` | Выход | ❌ |

### Athlete

| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| GET | `/athlete/profile` | Профиль | ✅ |
| PUT | `/athlete/profile` | Обновить профиль | ✅ |
| DELETE | `/athlete/profile` | Удалить аккаунт | ✅ |

### Exercises

| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| POST | `/exercise` | Создать упражнение | ✅ |
| GET | `/ListExercises?limit=10&page=1&filter=` | Список упражнений | ✅ |
| GET | `/exercise/{id}` | Упражнение по ID | ✅ |
| PUT | `/exercise/{id}` | Обновить упражнение | ✅ |
| DELETE | `/exercise/{id}` | Удалить упражнение | ✅ |

### Workouts (WIP)

| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| POST | `/workout` | Создать тренировку | ✅ |

## Примеры запросов

### Регистрация

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com","password":"secure123","name":"John","age":25}'
```

### Создание упражнения

```bash
curl -X POST http://localhost:8080/exercise \
  -H "Content-Type: application/json" \
  -b "access_token=<token>" \
  -d '{"name":"Bench Press","muscle_group":"Chest","description":"Classic exercise"}'
```

### Создание тренировки

```bash
curl -X POST http://localhost:8080/workout \
  -H "Content-Type: application/json" \
  -b "access_token=<token>" \
  -d '{
    "status": "finished",
    "total_time": 90,
    "grade_of_training": 8,
    "date_of_training": "2026-03-31T10:00:00Z",
    "sets": [
      {"exercise_id": 1, "weight": 100, "reps": 10, "set_order": 1, "rpe": 8},
      {"exercise_id": 1, "weight": 105, "reps": 8, "set_order": 2, "rpe": 9},
      {"exercise_id": 2, "weight": 80, "reps": 12, "set_order": 1, "rpe": 7}
    ]
  }'
```

## Тесты

```bash
go test ./...           # все тесты
go test ./... -v        # подробный вывод
go test ./... -cover    # покрытие
```

## Безопасность

- Пароли хешируются через **bcrypt** (cost factor 12)
- Access токен: 30 минут, Refresh токен: 30 дней
- Токены хранятся в **HttpOnly, Secure, SameSite=Strict** cookies
- Разделение access/refresh секретов
- Проверка прав на упражнения при создании тренировки

## Что в планах

- [ ] Получение тренировок по ID / пагинация
- [ ] Добавление сетов к `in_progress` тренировке
- [ ] Завершение тренировки
- [ ] Body measurements (замеры тела)
- [ ] Статистика и аналитика
- [ ] Swagger/OpenAPI документация
- [ ] Интеграционные тесты с testcontainers
