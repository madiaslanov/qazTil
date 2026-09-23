# qazTil API

HTTP API на Go для изучения казахского языка: категории, словарь, квиз и прогресс одного локального ученика. Авторизации нет. Страница живёт отдельно, в `../frontend`, и ходит сюда только за JSON.

`internal/domain` ни от кого не зависит, `internal/usecase` знает только домен, `internal/adapter` реализует HTTP (chi) и SQLite, `main.go` собирает их вместе.

Прод-хост: `https://qaztil.onrender.com/`.

## Запуск

Нужен Go 1.23 или новее. Из корня репозитория:

```bash
make run
```

Или из этого каталога: `go run .`.

Сервер слушает `http://localhost:8080`. Swagger UI: `http://localhost:8080/swagger/index.html`. Готовность: `GET /health`.

При пустой базе создаются три категории (приветствия, семья, еда) и слова с транскрипцией.

## Переменные окружения

- `ADDR` — адрес, по умолчанию `:8080`. Если задан `PORT` (так делает Render), сервер слушает `:$PORT`
- `DB_PATH` — файл SQLite, по умолчанию `data/qaztil.db`
- `ALLOWED_ORIGINS` — домены фронтенда через запятую, по умолчанию `http://localhost:3000`. На Render впишите сюда адрес Vercel
- `WEB_DIR` — каталог статики. По умолчанию пуст: сервер статику не отдаёт

## Связка с фронтендом

Клиент читает `NEXT_PUBLIC_API_URL` — базовый путь `/api/v1`, без завершающего слэша.

- локально (`next dev`) — `frontend/.env.development`: `http://localhost:8080/api/v1`
- прод (`next build`, Vercel) — `frontend/.env.production`: `https://qaztil.onrender.com/api/v1`

Шаблон обеих строк — в `frontend/.env.example`. Пока в `ALLOWED_ORIGINS` нет домена страницы, браузер ответы не примет.

## API

Базовый путь `/api/v1`. На проде тот же контракт: `https://qaztil.onrender.com/api/v1`.

- `GET/POST /categories`, `GET/PUT/DELETE /categories/{id}`
- `GET/POST /words?category_id=`, `GET/PUT/DELETE /words/{id}`
- `POST /quizzes` — тело `{"category_id": 1, "size": 5}`, в ответе русское слово и четыре казахских варианта
- `POST /quizzes/{id}/answers` — тело `{"question_id": 1, "selected_index": 0}`
- `GET /quizzes/{id}` — счёт
- `GET /progress` — ответы по каждой категории
- `GET /health` — готовность, вне `/api/v1`

Swagger на проде: `https://qaztil.onrender.com/swagger/index.html`.

## Деплой

Render, Root Directory — `backend`. Тот же контракт лежит в `../render.yaml`.

- Build Command: `go build -tags netgo -ldflags '-s -w' -o app`
- Start Command: `./app`
- Health Check Path: `/health`

Диск `/var/data` доступен на платном инстансе; на бесплатном уберите блок `disk` из `render.yaml`, тогда SQLite живёт только до следующего деплоя.

Чтобы чужой коммит фронтенда не пересобирал API, на Render в Build Filters включите путь `backend/**`.

## Docker

Образ собирается только из этого каталога и не включает страницу. Процесс идёт не от root, готовность проверяется через `GET /health`.

Из корня репозитория:

```bash
docker compose up --build
```

Контейнер готов, когда статус `healthy`.

## Проверки

Из корня репозитория:

```bash
make test
make build
make swagger
```
