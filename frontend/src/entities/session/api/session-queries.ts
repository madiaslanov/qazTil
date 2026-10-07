import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type {
  LoginPayload,
  RegisterPayload,
  Session,
  User,
} from "../model/types";

export const sessionQueries = {
  me: () =>
    queryOptions({
      queryKey: ["session", "me"],
      queryFn: () => request<User>("/auth/me"),
      // 401 повторять бессмысленно: клиент уже сбросил сессию.
      retry: false,
    }),
};

export const sessionApi = {
  login: (payload: LoginPayload) =>
    request<Session>("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  register: (payload: RegisterPayload) =>
    request<Session>("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};
