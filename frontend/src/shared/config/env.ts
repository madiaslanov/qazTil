/** Базовый адрес Go API. На Vercel задаётся в переменных окружения проекта. */
export const API_URL = (
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1"
).replace(/\/$/, "");

/** Swagger живёт рядом с API, на том же хосте. */
export const SWAGGER_URL = `${API_URL.replace(/\/api\/v1$/, "")}/swagger/index.html`;

/**
 * Пока эндпоинты курса не готовы, фронт работает на заглушках с теми же
 * формами, что в docs/api-contracts.md. Включается NEXT_PUBLIC_USE_MOCKS=true.
 */
export const USE_MOCKS = process.env.NEXT_PUBLIC_USE_MOCKS === "true";
