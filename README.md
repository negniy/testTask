```markdown
# People Management API

Это тестовое задание для Junior Golang Developer. Проект представляет собой REST API-сервис, который принимает ФИО, обогащает данные (возраст, пол, национальность) с помощью открытых API (agify.io, genderize.io, nationalize.io) и сохраняет информацию в базу данных PostgreSQL. Кроме того, сервис поддерживает операции поиска с фильтрами и пагинацией, обновление и удаление записей.

---

## Особенности

- **CRUD операции**: Добавление, получение, обновление и удаление записей о людях.
- **Обогащение данных**: При добавлении нового человека происходит вызов внешних API:
  - [agify.io](https://api.agify.io) – определение возраста по имени.
  - [genderize.io](https://api.genderize.io) – определение пола по имени.
  - [nationalize.io](https://api.nationalize.io) – определение национальности по имени.
- **Фильтрация и пагинация**: Возможность получать список людей с фильтрацией по полям (id, name, surname, patronymic, age, gender, nationality) и поддержкой пагинации (limit, offset).
- **Логирование**: Используется библиотека Logrus для ведения подробных логов (уровни debug и info). Логи записываются в файл `task.log`.
- **Конфигурация через .env**: Все настройки (например, параметры подключения к базе данных, порт сервера) хранятся в файле конфигурации `.env`.
- **Миграции**: При запуске приложения с помощью GORM автоматически выполняется миграция – создаются или обновляются таблицы в базе данных.
- **Swagger**: (Рекомендуется) можно сгенерировать документацию API с помощью swaggo для дальнейшей интеграции.

---

## Структура проекта

```
project/
├── cmd/
│   └── main.go         // Точка входа, настройка сервера, роутер и запуск HTTP-сервера
├── config/
│   └── config.go       // Настройка логгера и загрузка переменных окружения
├── handlers/
│   └── handlers.go     // Обработчики REST-запросов (GET, POST, PUT, DELETE)
├── internal/
│   └── enrich.go       // Логика вызова внешних API для обогащения данных
├── models/
│   └── models.go       // Модели данных (структура Person, структуры для обогащения)
├── repository/
│   └── repository.go   // Работа с базой данных (CRUD операции с использованием GORM)
├── .env                // Файл конфигурации (например, PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
└── task.log            // Файл логов
```

---

## Установка и настройка

### 1. Предварительные требования

- [Go](https://golang.org/dl/) (v1.16 или выше)
- [PostgreSQL](https://www.postgresql.org/download/) – база данных
- [Git](https://git-scm.com/)

### 2. Клонирование репозитория

```bash
git clone https://github.com/yourusername/people-management-api.git
cd people-management-api
```

### 3. Установка зависимостей

Используйте модули Go для установки зависимостей:

```bash
go mod tidy
```

### 4. Настройка переменных окружения

Создайте файл `.env` (или используйте уже существующий файл `config/conf.env`) в каталоге `config/` с примерно таким содержимым:

```
PORT=8080
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=task
DB_HOST=localhost
DB_PORT=5432
```

> **Важно:** GORM не создаёт базу данных автоматически – убедитесь, что база данных с именем `task` создана в PostgreSQL. Если база данных отсутствует, создайте её вручную или через скрипт:
> 
> ```sql
> CREATE DATABASE task;
> ```

---

## Запуск сервера

Для запуска сервера выполните:

```bash
go run main.go
```

Сервер запустится на порту, указанном в переменной окружения `PORT` (по умолчанию 8080). В логах вы увидите информацию о запуске и обработке запросов.

---

## API эндпоинты

### Получение списка людей (с фильтрами и пагинацией)

- **Метод:** GET  
- **URL:** `/people`  
- **Параметры запроса:**
  - `id` (необязательно) – фильтр по ID
  - `name` (необязательно) – фильтр по имени
  - `surname` (необязательно) – фильтр по фамилии
  - `patronymic` (необязательно) – фильтр по отчеству
  - `age` (необязательно) – фильтр по возрасту
  - `gender` (необязательно) – фильтр по полу
  - `nationality` (необязательно) – фильтр по национальности
  - `limit` (необязательно) – количество записей на страницу (по умолчанию 10)
  - `offset` (необязательно) – смещение для пагинации (по умолчанию 0)

**Пример запроса:**

```
GET /people?name=Dmitriy&age=30&limit=10&offset=0
```

### Создание нового человека

- **Метод:** POST  
- **URL:** `/people`  
- **Тело запроса (JSON):**

```json
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
}
```

При создании происходит обогащение данных (определение возраста, пола и национальности) с использованием внешних API.

### Обновление данных человека

- **Метод:** PUT  
- **URL:** `/people`  
- **Параметры запроса:** `id` (передается как параметр запроса, например: `/people?id=1`)  
- **Тело запроса (JSON):**

```json
{
  "name": "Новое имя",
  "surname": "Новая фамилия"
}
```

### Удаление человека

- **Метод:** DELETE  
- **URL:** `/people`  
- **Параметры запроса:** `id` (например: `/people?id=1`)

---

## Логирование

- Логирование осуществляется с помощью библиотеки **Logrus**.
- Логи записываются в файл `task.log` с уровнем Debug и выше.
- Логируется информация о входящих запросах, завершении их обработки, а также ошибки.

---

## Миграции базы данных

- При запуске приложения в функции `LoadDB()` происходит автоматическая миграция (с помощью метода `AutoMigrate` из GORM).  
- Если структура модели `Person` изменяется, GORM создаст или обновит соответствующую таблицу.

---

## Внутренняя логика обогащения данных

- Файл `internal/enrich.go` содержит функцию `EnrichPerson`, которая обращается к API [agify.io](https://api.agify.io), [genderize.io](https://api.genderize.io) и [nationalize.io](https://api.nationalize.io) для получения возраста, пола и национальности по имени.

---

## Технологии

- **Go** – язык программирования.
- **GORM** – ORM для работы с базой данных.
- **PostgreSQL** – СУБД.
- **Gorilla Mux** – роутер для обработки HTTP-запросов.
- **Logrus** – логирование.
- **godotenv** – загрузка переменных окружения.

---

## Заключение

Этот проект демонстрирует создание REST API-сервиса с обогащением данных, поддержкой фильтрации и пагинации, а также автоматическими миграциями базы данных. Проект легко настраивается через `.env` файл и включает подробное логирование для отладки и мониторинга.

---
