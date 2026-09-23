/** Базовый адрес Go API. На Vercel задаётся в переменных окружения проекта. */
export const API_URL = (
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1"
).replace(/\/$/, "");

/** Swagger живёт рядом с API, на том же хосте. */
export const SWAGGER_URL = `${API_URL.replace(/\/api\/v1$/, "")}/swagger/index.html`;
