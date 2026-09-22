# qazTil

Базовое приложение для изучения казахского языка. Словарь, категории, квиз и прогресс одного локального ученика. Авторизации нет.

Код сервера лежит в `backend/`: `internal/domain` ни от кого не зависит, `internal/usecase` знает только домен, `internal/adapter` реализует HTTP и SQLite, `cmd/api` собирает их вместе. Статичная страница — в `frontend/`.

## Запуск

Нужен Go 1.23 или новее.

```bash
make run
```

Сервер слушает `http://localhost:8080`. Страница-макет: `/`. Swagger UI: `http://localhost:8080/swagger/index.html`.

При пустой базе создаются три категории (приветствия, семья, еда) и слова с транскрипцией.

Переменные окружения:

- `ADDR` — адрес, по умолчанию `:8080`. Если задан `PORT` (так делает Render), сервер слушает `:$PORT`
- `DB_PATH` — файл SQLite, по умолчанию `data/qaztil.db`
- `WEB_DIR` — каталог статики, по умолчанию `../frontend` (относительно `backend/`)

## API

Базовый путь `/api/v1`.

- `GET/POST /categories`, `GET/PUT/DELETE /categories/{id}`
- `GET/POST /words?category_id=`, `GET/PUT/DELETE /words/{id}`
- `POST /quizzes` — тело `{"category_id": 1, "size": 5}`, в ответе русское слово и четыре казахских варианта
- `POST /quizzes/{id}/answers` — тело `{"question_id": 1, "selected_index": 0}`
- `GET /quizzes/{id}` — счёт
- `GET /progress` — ответы по каждой категории

Файлы в `frontend/` — статичный каркас. Логику интерфейса сюда подключит фронтенд.

## Docker

Образ в корне собирает API из `backend/` и кладёт рядом страницу из `frontend/`.

```bash
docker compose up --build
```

На Render в сервисе выберите среду Docker, а не Native Go. Корень репозитория — `.`, файл сборки — `Dockerfile`. База SQLite лежит в контейнере в `/app/data` и сбрасывается при новом деплое, если к сервису не подключён диск по этому пути.

## Проверки

```bash
make test
make build
make swagger
```
