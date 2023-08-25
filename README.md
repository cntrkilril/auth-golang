# Тестовое задание "auth-golang"

## Запуск приложения

В корне приложения в терминале в такой последовательности произвести команды:

1. `go mod tidy`
2. `make start-mongo`
3. `make run`

## Документация к api

1. GET /api/tokens/create-tokens/:user_id, где user_id - uuid пользователя
2. POST /api/tokens/refresh-tokens
   ``Body (JSON):
   {
"refreshToken": string,
"accessToken": string
   }
   ``