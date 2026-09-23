# qazTil

Базовое приложение для изучения казахского языка. Словарь, категории, квиз и прогресс одного локального ученика. Авторизации нет.

Код сервера лежит в `backend/`: `internal/domain` ни от кого не зависит, `internal/usecase` знает только домен, `internal/adapter` реализует HTTP и SQLite, `main.go` собирает их вместе. Клиент — в `frontend/`: Next.js (App Router) по архитектуре Feature-Sliced Design.

Фронтенд и бэкенд деплоятся отдельно: страница — на Vercel, API — на Render. Поэтому сервер отдаёт только JSON и разрешает запросы с домена фронтенда через CORS.

## Запуск

Нужен Go 1.23 или новее.

```bash
make run
```

Сервер слушает `http://localhost:8080`. Swagger UI: `http://localhost:8080/swagger/index.html`. Страница открывается отдельно, из `frontend/`.

При пустой базе создаются три категории (приветствия, семья, еда) и слова с транскрипцией.

Переменные окружения:

- `ADDR` — адрес, по умолчанию `:8080`. Если задан `PORT` (так делает Render), сервер слушает `:$PORT`
- `DB_PATH` — файл SQLite, по умолчанию `data/qaztil.db`
- `ALLOWED_ORIGINS` — домены фронтенда через запятую, по умолчанию `http://localhost:3000`. На Render впишите сюда адрес Vercel
- `WEB_DIR` — каталог статики. По умолчанию пуст: сервер статику не отдаёт

## API

Базовый путь `/api/v1`.

- `GET/POST /categories`, `GET/PUT/DELETE /categories/{id}`
- `GET/POST /words?category_id=`, `GET/PUT/DELETE /words/{id}`
- `POST /quizzes` — тело `{"category_id": 1, "size": 5}`, в ответе русское слово и четыре казахских варианта
- `POST /quizzes/{id}/answers` — тело `{"question_id": 1, "selected_index": 0}`
- `GET /quizzes/{id}` — счёт
- `GET /progress` — ответы по каждой категории

## Фронтенд

Next.js 16 (App Router), Tailwind v4, shadcn/ui для базовых примитивов, TanStack Query для запросов к API и Zustand для локального профиля ученика. Вёрстка собрана по макету Figma «QazTil».

Слои FSD лежат в `frontend/src/`:

- `app/` — роутер Next и провайдеры
- `views/` — экраны (слой страниц FSD; имя `pages` занято роутером)
- `widgets/` — шапка, нижняя навигация, путь обучения
- `features/` — онбординг и прохождение урока
- `entities/` — категории, слова, квиз, прогресс, профиль ученика
- `shared/` — HTTP-клиент, токены дизайна, UI-кит

Экраны: `/` — онбординг, `/learn` — путь обучения, `/lesson/{categoryId}` — урок и его итог, `/progress` — прогресс по категориям, `/profile` — профиль, `/words` — словарь.

Авторизации, страйков, жизней и XP в API нет, поэтому профиль ученика целиком живёт в браузере (`localStorage`), а пароль с экрана входа никуда не отправляется.

### Деплой

Сервисы не собирают файлы друг друга.

**Render, API.** Root Directory — `backend`. Команду сборки оставьте ту, что Render ставит сам: `go build -tags netgo -ldflags '-s -w' -o app`. Start Command: `./app`. Health Check Path: `/health`. В `ALLOWED_ORIGINS` впишите адрес Vercel. Тот же контракт лежит в `render.yaml`. Диск `/var/data` доступен на платном инстансе; на бесплатном уберите блок `disk`, тогда SQLite живёт только до следующего деплоя.

**Vercel, страница.** Root Directory — `frontend/`, сборка стандартная (`next build`). В переменных окружения задайте `NEXT_PUBLIC_API_URL`, например `https://qaztil.onrender.com/api/v1`.

Чтобы чужой коммит не пересобирал сервис: на Render в Build Filters включите путь `backend/**`, на Vercel в Ignored Build Step — `git diff HEAD^ HEAD --quiet -- frontend`.

## Docker

Образ собирается только из `backend/` и не включает страницу. Процесс идёт не от root, готовность проверяется через `GET /health`.

```bash
docker compose up --build
docker compose ps
```

Контейнер готов, когда статус `healthy`. Swagger: `http://localhost:8080/swagger/index.html`.

## Проверки

```bash
make test
make build
make swagger
make web-build
make web-lint
```
