"use client";

import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";

import { useLearnerStore, type LearnerPrefs } from "@/entities/learner";
import {
  sessionApi,
  useSessionStore,
  type LoginPayload,
} from "@/entities/session";
import { ApiError } from "@/shared/api";

type AuthRequest =
  | { mode: "login"; credentials: LoginPayload }
  | { mode: "register"; credentials: LoginPayload; prefs: LearnerPrefs };

/** Вход и регистрация: сохраняем сессию, заводим профиль, ведём на путь. */
export function useAuth() {
  const router = useRouter();
  const setSession = useSessionStore((state) => state.setSession);
  const learner = useLearnerStore((state) => state.learner);
  const createLearner = useLearnerStore((state) => state.create);

  const mutation = useMutation({
    // Имя не спрашиваем — Go возьмёт его из почты.
    mutationFn: (request: AuthRequest) =>
      request.mode === "register"
        ? sessionApi.register({ ...request.credentials, display_name: "" })
        : sessionApi.login(request.credentials),
    onSuccess: (session, request) => {
      setSession(session);
      // Геймификация живёт в браузере: новый профиль заводим при регистрации
      // и при входе с чужого для этого устройства аккаунта.
      if (request.mode === "register") {
        createLearner(session.user.email, request.prefs);
      } else if (learner?.email !== session.user.email) {
        createLearner(session.user.email, { dailyGoal: 10 });
      }
      router.push("/learn");
    },
  });

  const error =
    mutation.error instanceof ApiError
      ? mutation.error.message
      : mutation.error && "Не получилось связаться с сервером";

  return { ...mutation, error };
}
