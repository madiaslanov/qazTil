# qazTil API

HTTP API на Go для изучения казахского языка: категории, словарь, ситуативные слова, юниты с уроками и упражнениями, квиз и прогресс одного локального ученика. Авторизации нет. Страница живёт отдельно, в `../frontend`, и ходит сюда только за JSON.

`internal/domain` ни от кого не зависит, `internal/usecase` знает только домен, `internal/adapter` реализует HTTP (chi) и SQLite, `main.go` собирает их вместе.

Прод-хост: `https://qaztil.onrender.com/`.

## Запуск

Нужен Go 1.23 или новее. Из корня репозитория:

```bash
make run
```

Или из этого каталога: `go run .`.

Сервер слушает `http://localhost:8080`. Swagger UI: `http://localhost:8080/swagger/index.html`. Готовность: `GET /health`.

При пустой базе создаются три категории (приветствия, семья, еда) и слова с транскрипцией. Отдельно, если ещё нет ситуаций, создаются три ситуации (кафе, магазин, дорога), слова к ним и юнит с двумя уроками на каждую. То же самое делает `make seed` из корня (`go run ./cmd/seed` из этого каталога): уже заполненные таблицы он не трогает.

## Переменные окружения

- `ADDR` — адрес, по умолчанию `:8080`. Если задан `PORT` (так делает Render), сервер слушает `:$PORT`
- `DB_PATH` — файл SQLite, по умолчанию `data/qaztil.db`
- `ALLOWED_ORIGINS` — домены фронтенда через запятую. По умолчанию `http://localhost:3000` и `https://qaz-til.vercel.app`
- `WEB_DIR` — каталог статики. По умолчанию пуст: сервер статику не отдаёт

## Связка с фронтендом

Клиент читает `NEXT_PUBLIC_API_URL` — базовый путь `/api/v1`, без завершающего слэша.

- локально (`next dev`) — `frontend/.env.development`: `http://localhost:8080/api/v1`
- прод (`next build`, Vercel) — `frontend/.env.production`: `https://qaztil.onrender.com/api/v1`

Шаблон обеих строк — в `frontend/.env.example`. Страница на проде — `https://qaz-til.vercel.app`, этот origin уже входит в `ALLOWED_ORIGINS`.

## API

Базовый путь `/api/v1`. На проде тот же контракт: `https://qaztil.onrender.com/api/v1`.

- `GET/POST /categories`, `GET/PUT/DELETE /categories/{id}`
- `GET/POST /words?category_id=`, `GET/PUT/DELETE /words/{id}`
- `POST /quizzes` — тело `{"category_id": 1, "size": 5}`, в ответе русское слово и четыре казахских варианта
- `POST /quizzes/{id}/answers` — тело `{"question_id": 1, "selected_index": 0}`
- `GET /quizzes/{id}` — счёт
- `GET /progress` — ответы по каждой категории
- `GET /situations`, `GET /situations/{id}` — ситуация и её слова; `GET /situations/{id}/words` — только слова
- `GET /units`, `GET /units/{id}` — юнит и уроки; `GET /units/{id}/lessons` — только уроки
- `GET /lessons/{id}` — урок и упражнения без правильного ответа; `GET /lessons/{id}/exercises` — только упражнения
- `POST /exercises/{id}/answers` — тело `{"answer": "Рақмет"}`, в ответе `correct` и ожидаемая строка
- `GET /health` — готовность, вне `/api/v1`

Упражнение бывает `choice` (казахские варианты на русский prompt) или `translate` (нужно написать казахскую строку). Юнит ссылается на ситуацию: слова ситуации и уроки юнита — один сценарий.

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
