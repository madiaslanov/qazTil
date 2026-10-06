#!/usr/bin/env bash
# Дымовой тест аутентификации QazTil (QAZ-22).
#
# Запуск:
#   export JWT_SECRET=dev-secret-change-me-0123456789abcdef
#   make run            # в другой вкладке
#   ./smoke-auth.sh
#
# Адрес можно переопределить:
#   BASE_URL=https://qaztil.onrender.com ./smoke-auth.sh
#
# Тела запросов собираются в переменные отдельными строками, а не внутри
# подстановки команд: bash 3.2 (штатный на macOS) рвёт такую вложенность
# по запятым и отправляет мусор.

set -u

BASE_URL="${BASE_URL:-http://localhost:8080}"
API="$BASE_URL/api/v1"

STAMP=$(date +%s)
EMAIL="smoke-$STAMP@example.com"
PASSWORD="qazaqtili2026"

failures=0

login_json() {
	printf '{"email":"%s","password":"%s"}' "$1" "$2"
}

register_json() {
	printf '{"email":"%s","password":"%s","display_name":"%s"}' "$1" "$2" "$3"
}

post_code() {
	curl -s -o /dev/null -w '%{http_code}' -X POST "$1" \
		-H 'Content-Type: application/json' -d "$2"
}

get_code() {
	curl -s -o /dev/null -w '%{http_code}' "$@"
}

# check <описание> <ожидаемый код> <фактический код>
check() {
	if [ "$2" = "$3" ]; then
		printf '  %-26s %s\n' "$1" "$3"
	else
		printf '  %-26s %s  <-- ОЖИДАЛОСЬ %s\n' "$1" "$3" "$2"
		failures=$((failures + 1))
	fi
}

echo "QazTil auth smoke test"
echo "  адрес: $BASE_URL"
echo "  почта: $EMAIL"
echo

if ! curl -sf -o /dev/null "$BASE_URL/health"; then
	echo "Сервер недоступен на $BASE_URL — запусти 'make run' и повтори."
	exit 1
fi

echo "Коды ответов:"

body=$(register_json "$EMAIL" "$PASSWORD" "Smoke")
code=$(post_code "$API/auth/register" "$body")
check "register new user" 201 "$code"

body=$(login_json "$EMAIL" "$PASSWORD")
code=$(post_code "$API/auth/login" "$body")
check "login correct password" 200 "$code"

login_response=$(curl -s -X POST "$API/auth/login" \
	-H 'Content-Type: application/json' -d "$body")
token=$(printf '%s' "$login_response" |
	sed -n 's/.*"token"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')

if [ -z "$token" ]; then
	echo "  не удалось извлечь токен из ответа логина:"
	echo "    $login_response"
	failures=$((failures + 1))
fi

code=$(get_code "$API/auth/me" -H "Authorization: Bearer $token")
check "me with token" 200 "$code"

code=$(get_code "$API/auth/me")
check "me without token" 401 "$code"

code=$(get_code "$API/auth/me" -H "Authorization: Bearer not.a.token")
check "me with broken token" 401 "$code"

body=$(login_json "$EMAIL" "wrongpassword")
code=$(post_code "$API/auth/login" "$body")
check "login wrong password" 401 "$code"

body=$(login_json "nobody-$STAMP@example.com" "$PASSWORD")
code=$(post_code "$API/auth/login" "$body")
check "login unknown email" 401 "$code"

# Тот же адрес в верхнем регистре: проверяем нормализацию почты.
upper_email=$(printf '%s' "$EMAIL" | tr '[:lower:]' '[:upper:]')
body=$(register_json "$upper_email" "$PASSWORD" "Smoke")
code=$(post_code "$API/auth/register" "$body")
check "register duplicate email" 409 "$code"

body=$(register_json "short-$STAMP@example.com" "qazaq" "Smoke")
code=$(post_code "$API/auth/register" "$body")
check "register short password" 400 "$code"

echo
echo "Утечка хеша пароля:"
if printf '%s' "$login_response" | grep -q 'password_hash'; then
	echo "  password_hash В ОТВЕТЕ  <-- так быть не должно"
	failures=$((failures + 1))
else
	echo "  password_hash отсутствует в ответе"
fi

echo
if [ "$failures" -eq 0 ]; then
	echo "Все проверки пройдены."
	exit 0
fi
echo "Провалено проверок: $failures"
exit 1