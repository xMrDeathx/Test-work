# Тестовое задание "BackDev"

## Задание

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

## Требования

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

## Для ручного запуска:
1. `go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0`
2. `oapi-codegen -generate gorilla,types authentication/api/frontend/frontendapi.yaml > authentication/api/frontend/frontendapi.gen.go`
3. `go mod tidy -v`
4. `go run main.go`

## Для работы с БД
1. Установить PostgreSQL
2. Миграции выполняются автоматически

## Миграции
up и down миграции выполняются автоматически при запуске приложения

[Инструкция по установке утилиты для миграций](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)

Создание миграции: `migrate create -ext sql -dir {path-to-migrations-dir} {migartion-name}`

## Основные маршруты API

- **GET /authentication/login** - Получение access и refresh токенов.
- **GET /authentication/refresh-token/{userId}** - Обновление токенов по ID пользователя.

## Зависимости

Проект использует следующие библиотеки:

- **Gorilla Mux**: HTTP-маршрутизатор для Go.
- **pgx**: Драйвер для работы с базой данных PostgreSQL.
- **oapi**: Инструмент для генерации Go-кода на основе спецификаций OpenAPI.
- **godotenv**: Для загрузки переменных окружения из файла .env.
