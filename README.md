# People Management API

Этот проект представляет собой REST API-сервис для управления данными о людях. Сервис принимает ФИО, обогащает данные (возраст, пол, национальность) с помощью открытых API (agify.io, genderize.io, nationalize.io) и сохраняет информацию в базу данных PostgreSQL. Кроме того, сервис поддерживает операции поиска с фильтрами и пагинацией, обновление и удаление записей.

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
- **Swagger документация**: Сгенерированная документация API с помощью swaggo (swag). Swagger UI доступен для просмотра и тестирования API.

---

## Структура проекта

```
project/
├── cmd/
│   └── main.go         // Точка входа, настройка сервера, роутер и запуск HTTP-сервера.
├── config/
│   └── config.go       // Настройка логгера и загрузка переменных окружения.
├── handlers/
│   └── handlers.go     // Обработчики REST-запросов (GET, POST, PUT, DELETE).
├── internal/
│   └── enrich.go       // Логика вызова внешних API для обогащения данных.
├── models/
│   └── models.go       // Модели данных (структура Person, структуры для обогащения).
├── repository/
│   └── repository.go   // Работа с базой данных (CRUD операции с использованием GORM).
├── config/conf.env     // Файл конфигурации (например, PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT).
└── task.log            // Файл логов.
```

---

## Установка и настройка

### Настройка переменных окружения

Создайте файл `.env` (или используйте уже существующий файл `config/conf.env`) в каталоге `config/` со следующим содержимым:

```
PORT=8080
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=task
DB_HOST=localhost
DB_PORT=5432
```

> **Важно:** Убедитесь, что база данных с именем `task` создана в PostgreSQL. Если база данных отсутствует, создайте её вручную или с помощью скрипта:
> ```sql
> CREATE DATABASE task;
> ```

## Запуск приложения

Для запуска сервера выполните:

```bash
go run main.go
```

Сервер запустится на порту, указанном в переменной окружения `PORT` (по умолчанию 8080). В логах будут отображаться сообщения о запуске и обработке запросов, а Swagger UI будет доступен для просмотра документации по адресу `localhost:PORT\swagger\`.

---

## API эндпоинты

### Получение списка людей

- **Метод:** GET  
- **URL:** `/people`  
- **Параметры запроса (необязательно):**
  - `id` — фильтр по ID
  - `name` — фильтр по имени
  - `surname` — фильтр по фамилии
  - `patronymic` — фильтр по отчеству
  - `age` — фильтр по возрасту
  - `gender` — фильтр по полу
  - `nationality` — фильтр по национальности
  - `limit` — число записей на странице (по умолчанию 10)
  - `offset` — смещение для пагинации (по умолчанию 0)

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

При создании происходит обогащение данных (возраст, пол, национальность) с использованием внешних API.

### Обновление данных человека

- **Метод:** PUT  
- **URL:** `/people` (ID передается как query-параметр, например: `/people?id=1`)  
- **Тело запроса (JSON):**

```json
{
  "name": "Новое имя",
  "surname": "Новая фамилия"
}
```

### Удаление человека

- **Метод:** DELETE  
- **URL:** `/people` (ID передается как query-параметр, например: `/people?id=1`)

---

## Логирование

- Используется библиотека **Logrus**.
- Логи записываются в файл `task.log` с уровнем Debug и выше.
- Логируется информация о входящих запросах, завершении их обработки и возникающих ошибках.

---

## Миграции базы данных

При запуске приложения вызывается функция `LoadDB()`, которая устанавливает соединение с PostgreSQL и автоматически выполняет миграцию с помощью GORM (через `db.AutoMigrate(&models.Person{})`). Это создает или обновляет таблицу `people` в базе данных.

---

## Обогащение данных

Функция `EnrichPerson` в файле `internal/enrich.go` обращается к следующим API для получения дополнительных данных:
- **agify.io** – определение возраста.
- **genderize.io** – определение пола.
- **nationalize.io** – определение национальности.

Эти данные добавляются к создаваемым записям о людях.

---

## Заключение

Проект демонстрирует создание REST API-сервиса с обогащением данных, поддержкой фильтрации и пагинации, автоматическими миграциями базы данных и интеграцией Swagger-документации. Проект легко настраивается через файл `.env`, а подробное логирование помогает в отладке и мониторинге работы сервера.

---