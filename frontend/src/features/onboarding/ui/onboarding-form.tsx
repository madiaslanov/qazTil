"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";

import { useLearnerStore, type DailyGoal } from "@/entities/learner";
import { sessionApi, useSessionStore } from "@/entities/session";
import { ApiError } from "@/shared/api";
import { Button, Input, Label } from "@/shared/ui";

import { GoalPicker } from "./goal-picker";

type Mode = "register" | "login";

const MIN_PASSWORD = 8;
/** bcrypt на сервере режет всё после 72 байт, Go такой пароль не примет. */
const MAX_PASSWORD_BYTES = 72;

function passwordError(password: string): string | null {
  if (password.length < MIN_PASSWORD) {
    return `пароль — минимум ${MIN_PASSWORD} символов`;
  }
  if (new TextEncoder().encode(password).length > MAX_PASSWORD_BYTES) {
    return "пароль слишком длинный";
  }
  return null;
}

/**
 * Вход из макета: регистрация с дневной целью или вход в существующий аккаунт.
 * Имя не спрашиваем — Go возьмёт его из почты.
 */
export function OnboardingForm() {
  const router = useRouter();
  const learner = useLearnerStore((state) => state.learner);
  const createLearner = useLearnerStore((state) => state.create);
  const setSession = useSessionStore((state) => state.setSession);
  const [mode, setMode] = useState<Mode>("register");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [goal, setGoal] = useState<DailyGoal>(10);
  const [localError, setLocalError] = useState<string | null>(null);

  const auth = useMutation({
    mutationFn: () => {
      const credentials = { email: email.trim(), password };
      return mode === "register"
        ? sessionApi.register({ ...credentials, display_name: "" })
        : sessionApi.login(credentials);
    },
    onSuccess: (session) => {
      setSession(session);
      // Геймификация пока живёт в браузере: заводим профиль под этот аккаунт,
      // а прогресс того же ученика на этом устройстве не трогаем.
      if (mode === "register" || learner?.email !== session.user.email) {
        createLearner(session.user.email, goal);
      }
      router.push("/learn");
    },
  });

  function submit(event: React.FormEvent) {
    event.preventDefault();
    const error = passwordError(password);
    setLocalError(error);
    if (!error) auth.mutate();
  }

  function switchMode() {
    setMode(mode === "register" ? "login" : "register");
    setLocalError(null);
    auth.reset();
  }

  const error =
    localError ??
    (auth.error instanceof ApiError
      ? auth.error.message
      : auth.error && "ошибка запроса");

  return (
    <form onSubmit={submit} className="flex flex-1 flex-col gap-6.5">
      <div className="flex flex-col gap-3.5">
        <div className="flex flex-col gap-1.75">
          <Label htmlFor="email">Почта</Label>
          <Input
            id="email"
            type="email"
            required
            autoComplete="email"
            placeholder="you@example.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
        </div>
        <div className="flex flex-col gap-1.75">
          <Label htmlFor="password">Пароль</Label>
          <Input
            id="password"
            type="password"
            required
            autoComplete={
              mode === "register" ? "new-password" : "current-password"
            }
            placeholder="••••••••"
            aria-invalid={Boolean(error)}
            aria-describedby={error ? "auth-error" : undefined}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </div>
        {error && (
          <p
            id="auth-error"
            role="alert"
            className="text-caption font-bold text-accent"
          >
            {error}
          </p>
        )}
      </div>

      {mode === "register" && <GoalPicker value={goal} onChange={setGoal} />}

      <div className="mt-auto flex flex-col gap-3 pt-6">
        <Button type="submit" variant="primary" disabled={auth.isPending}>
          {mode === "register" ? "Начать обучение" : "Войти"}
        </Button>
        <button
          type="button"
          onClick={switchMode}
          className="text-body-sm font-bold text-muted underline underline-offset-4"
        >
          {mode === "register"
            ? "Уже есть аккаунт? Войти"
            : "Нет аккаунта? Зарегистрироваться"}
        </button>
      </div>
    </form>
  );
}
