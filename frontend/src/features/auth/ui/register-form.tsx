"use client";

import { useState } from "react";
import Link from "next/link";
import { CircleAlert } from "lucide-react";

import type { LearnerPrefs } from "@/entities/learner";
import { Button, FieldHint, FormField, Input } from "@/shared/ui";

import { emailError, passwordError } from "../lib/validation";
import { useAuth } from "../model/use-auth";

/**
 * Регистрация: почта и пароль. Настройки ученика собирает онбординг,
 * его поля встают между паролем и кнопкой через children.
 */
export function RegisterForm({
  prefs,
  children,
}: {
  prefs: LearnerPrefs;
  children?: React.ReactNode;
}) {
  const auth = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailTouched, setEmailTouched] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);

  const emailProblem = emailTouched ? emailError(email) : null;
  const error = localError ?? auth.error;

  function submit(event: React.FormEvent) {
    event.preventDefault();
    setEmailTouched(true);
    const problem = emailError(email) ?? passwordError(password);
    setLocalError(emailError(email) ? null : problem);
    if (problem) return;
    auth.mutate({
      mode: "register",
      credentials: { email: email.trim(), password },
      prefs,
    });
  }

  return (
    <form onSubmit={submit} noValidate className="flex flex-1 flex-col">
      <div className="flex flex-col gap-3.5">
        <FormField
          id="email"
          label="Почта"
          hint={emailProblem}
          hintTone="error"
        >
          <Input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            value={email}
            status={emailProblem ? "error" : "default"}
            icon={emailProblem && <CircleAlert />}
            aria-describedby={emailProblem ? "email-hint" : undefined}
            onChange={(event) => setEmail(event.target.value)}
            onBlur={() => setEmailTouched(true)}
          />
        </FormField>

        <FormField id="password" label="Пароль">
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            placeholder="••••••••"
            value={password}
            status={localError ? "error" : "default"}
            aria-describedby={error ? "auth-error" : undefined}
            onChange={(event) => setPassword(event.target.value)}
          />
        </FormField>

        {error && (
          <FieldHint id="auth-error" tone="error">
            {error}
          </FieldHint>
        )}
      </div>

      {children && <div className="mt-6.5">{children}</div>}

      <div className="mt-auto flex flex-col items-center gap-4 pt-8">
        <Button type="submit" disabled={auth.isPending}>
          {auth.isPending ? "Создаём профиль…" : "Начать обучение"}
        </Button>
        <Button asChild variant="ghost" size="sm">
          <Link href="/login">У меня уже есть аккаунт</Link>
        </Button>
      </div>
    </form>
  );
}
