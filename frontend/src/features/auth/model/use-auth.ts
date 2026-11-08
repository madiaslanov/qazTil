"use client";

import { useRouter } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import {
  profileQueries,
  useLearnerStore,
  type LearnerPrefs,
} from "@/entities/learner";
import {
  sessionApi,
  useSessionStore,
  type LoginPayload,
} from "@/entities/session";
import { ApiError } from "@/shared/api";

type AuthRequest =
  | { mode: "login"; credentials: LoginPayload }
  | { mode: "register"; credentials: LoginPayload; prefs: LearnerPrefs };

/** Вход и регистрация: сессия в localStorage, prefs — на сервер. */
export function useAuth() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const setSession = useSessionStore((state) => state.setSession);
  const learner = useLearnerStore((state) => state.learner);
  const createLearner = useLearnerStore((state) => state.create);

  const mutation = useMutation({
    mutationFn: (request: AuthRequest) =>
      request.mode === "register"
        ? sessionApi.register({
            ...request.credentials,
            display_name: "",
            daily_goal: request.prefs.dailyGoal,
            level: request.prefs.level,
            motivation: request.prefs.motivation,
          })
        : sessionApi.login(request.credentials),
    onSuccess: (session, request) => {
      setSession(session);
      // Локально — только прогресс по категориям; prefs уже в Go.
      if (
        request.mode === "register" ||
        learner?.email !== session.user.email
      ) {
        createLearner(session.user.email);
      }
      void queryClient.invalidateQueries({
        queryKey: profileQueries.current().queryKey,
      });
      router.push("/learn");
    },
  });

  const error =
    mutation.error instanceof ApiError
      ? mutation.error.message
      : mutation.error && "Не получилось связаться с сервером";

  return { ...mutation, error };
}
