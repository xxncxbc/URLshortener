# 📚 API URL Shortener

## Общая информация
- **База URL API**: `http://localhost:8080`
- **Аутентификация**: JWT Access Token and Refresh Token
- **Основные пакеты**: net/http, gorm, bcrypt,
  golang-jwt/jwt/v5, joho/godotenv, DATA-DOG/go-sqlmock,
  go-playground/validator/v10

## 🧑‍💻 Аутентификация

### POST `/auth/login`
Авторизация пользователя и получение токенов.

**Запрос:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "access_token": "JWT_ACCESS_TOKEN",
  "refresh_token": "JWT_REFRESH_TOKEN"
}
```

### POST `/auth/register`
Регистрация пользователя и получение токенов.

**Запрос:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Ответ:**
```json
{
  "access_token": "JWT_ACCESS_TOKEN",
  "refresh_token": "JWT_REFRESH_TOKEN"
}
```

### POST `/auth/refresh`
Обновление токенов по refresh token.

**Запрос:**
```json
{
  "refresh_token": "JWT_REFRESH_TOKEN"
}
```

**Ответ:**
```json
{
  "access_token": "NEW_ACCESS_TOKEN",
  "refresh_token": "NEW_REFRESH_TOKEN"
}
```

## 🔗 Работа с короткими ссылками

### POST `/link`
Создать короткий URL.

**Заголовки:**
- Authorization: `{access_token}`

**Запрос:**
```json
{
  "url": "https://example.com/very/long/url"
}
```

**Ответ:**
```json
{
    "ID": 10,
    "CreatedAt": "2025-03-17T00:50:50.5987413+03:00",
    "UpdatedAt": "2025-03-17T00:50:50.5987413+03:00",
    "DeletedAt": null,
    "url": "https://photomath.com",
    "hash": "xSxInejPNy",
    "user_id": 7
}
```
---

### GET `/{hash}`
Переход по короткой ссылке на оригинальный URL.

**Пример:**  
`GET /abc123`  
**Ответ:**  
`302 Redirect` на оригинальный URL.

---

### GET `/link`
Получить все ссылки

**Заголовки:**
- Authorization: `{access_token}`

- **Ответ:**
```json
{
    "ID": 10,
    "CreatedAt": "2025-03-17T00:50:50.5987413+03:00",
    "UpdatedAt": "2025-03-17T00:50:50.5987413+03:00",
    "DeletedAt": null,
    "url": "https://photomath.com",
    "hash": "xSxInejPNy",
    "user_id": 7
}
```

---

### DELETE `/link/{id}`
Удалить ссылку по ID.

**Заголовки:**
- Authorization:`{access_token}`

**Ответ:**
- `204 No Content` или ошибка (403, 500)

---

### PATCH `/link/{id}`

**Заголовки:**
- Authorization:`{access_token}`

**Запрос:**
```json
{
    "url": "https://google.com",
    "hash": "GGGGG"
}
```
**Ответ:**
```json
{
    "ID": 3,
    "CreatedAt": "2025-03-17T01:15:28.628897+03:00",
    "UpdatedAt": "2025-03-17T01:16:33.415534+03:00",
    "DeletedAt": null,
    "url": "https://google.com",
    "hash": "GGGGG",
    "user_id": 1
}
```
## 📊 Статистика

### GET `/stat?from=(YYYY-MM-DD/YYYY-MM)&to=(YYYY-MM-DD/YYYY-MM)&by=month`
Общая статистика по всем ссылкам, где sum это количество переходов по ссылкам.

**Заголовки:**
- Authorization:`{access_token}`

**Ответ:**
```json
[
    {
        "period": "2025-03",
        "sum": 12
    }
]
```

---

### ⚙️ Middleware в проекте

Проект использует собственные middleware для централизованной обработки HTTP-запросов. Они упрощают логику и позволяют использовать кросс-срезовую функциональность (например, логирование или аутентификацию) для всех или отдельных роутов.

#### Основные middleware:

---

### 1️⃣ **Auth Middleware (`IsAuthed`)**

Используется для проверки JWT и извлечения данных о пользователе (email и userId) из токена. Если токен невалиден — возвращается `401 Unauthorized`. При успешной валидации данные пользователя прокидываются в контекст запроса.

```go
func IsAuthed(next http.Handler, config *configs.Config) http.Handler
```

Доступ к данным из контекста:
```go
email := r.Context().Value(middleware.ContextEmailKey).(string)
userId := r.Context().Value(middleware.ContextUserIdKey).(uint)
```

---

### 2️⃣ **CORS Middleware**

Добавляет необходимые заголовки для поддержки CORS. Также корректно обрабатывает `OPTIONS` preflight-запросы.

```go
func CORS(next http.Handler) http.Handler
```

Разрешенные методы:
- `GET, POST, PUT, DELETE, HEAD, PATCH`

Разрешенные заголовки:
- `authorization, content-type, content-length`

---

### 3️⃣ **Logging Middleware**

Логирует:
- HTTP метод
- Путь запроса
- HTTP статус код
- Время обработки запроса

Пример логов:
```
200 GET /api/links 5.231ms
```

```go
func Logging(next http.Handler) http.Handler
```

Используется обертка для ResponseWriter — `WrapperWriter` для корректного получения HTTP статуса.

---

### 4️⃣ **Chain Middleware**

Функция объединения нескольких middleware в единую цепочку:

```go
stack := middleware.Chain(
	middleware.CORS,
	middleware.Logging,
)
```

Вызов middleware будет происходить в обратном порядке, т.е. сначала `Logging`, затем `CORS`, и только потом передача запроса в handler.

---

### ✅ Пример использования всех middleware

```go
mux := http.NewServeMux()
mux.Handle("/api/protected", middleware.IsAuthed(yourHandler, config))

wrapped := middleware.Chain(
	middleware.CORS,
	middleware.Logging,
)(mux)

http.ListenAndServe(":8080", wrapped)
```
---


### ⚙ Асинхронный подсчет кликов с использованием шины событий

В системе реализован асинхронный подсчет кликов с помощью собственного event bus:

- При переходе по сокращённой ссылке сервис публикует событие `link_visited` в шину событий.
- В фоновом режиме работает отдельная горутина, которая подписана на события шины и обрабатывает их без блокировки основного HTTP-сервера.
- Это позволяет быстро отдавать клиенту ответ, а запись данных о клике происходит параллельно, не задерживая обработку запроса.

#### Основные компоненты:

- **EventBus** – собственная шина событий с методами `Publish` и `Subscribe`.
- **StatService** – сервис, который подписывается на события и асинхронно увеличивает количество кликов в базе данных через `StatRepository`.

Такой подход улучшает производительность и отзывчивость API при большом количестве трафика.

--- 


### 🧪 Тестирование

В проекте используется несколько видов тестов для обеспечения надёжности и корректности работы API:

#### End-to-End (E2E) тесты

Полноценные сквозные тесты, которые проверяют работу API вместе с базой данных и сервером.

Пример:

- Проверка успешной авторизации с реальной базой данных и HTTP-запросом на `/auth/login`.
- Проверка случая с неправильным паролем.

Файл: `auth_test.go`

---

#### Unit-тесты

Локальные тесты отдельных компонентов приложения без взаимодействия с внешними сервисами.

Пример:

- Тестирование логики создания и парсинга JWT-токена.

Файл: `jwthelper/jwt_test.go`

---

#### Mock-тесты

Используются мок-объекты и `sqlmock` для имитации работы с базой данных без реального подключения.

Примеры:

- Проверка успешного логина и регистрации пользователя с подменой реального репозитория на mock.
- Проверка сервиса регистрации с кастомным мок-репозиторием.

Файлы:
- `auth/handler_test.go`
- `auth/service_test.go`

---
