import { API_URL } from "@/shared/config/env";

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

type AuthBridge = {
  getToken: () => string | null;
  onUnauthorized: () => void;
};

/**
 * shared не знает про сессию: токен и реакцию на 401
 * подключает слой app через configureAuth.
 */
let auth: AuthBridge = {
  getToken: () => null,
  onUnauthorized: () => {},
};

export function configureAuth(bridge: AuthBridge) {
  auth = bridge;
}

/** Единственная точка запросов к Go API: JSON туда, JSON обратно. */
export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const token = auth.getToken();
  const headers = new Headers(init?.headers);
  if (init?.body) headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  let response: Response;
  try {
    response = await fetch(`${API_URL}${path}`, {
      ...init,
      headers,
      cache: "no-store",
    });
  } catch {
    throw new ApiError("сервер недоступен", 0);
  }

  if (!response.ok) {
    // 401 без токена — это неверный пароль при входе, а не протухшая сессия.
    if (response.status === 401 && token) {
      auth.onUnauthorized();
    }
    const body = await response.json().catch(() => null);
    const message =
      body && typeof body.error === "string" ? body.error : "ошибка запроса";
    throw new ApiError(message, response.status);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}
