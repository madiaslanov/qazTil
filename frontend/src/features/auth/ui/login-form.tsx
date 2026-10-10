"use client";

import { useState } from "react";
import Link from "next/link";
import { CircleAlert } from "lucide-react";

import { Button, FieldHint, FormField, Input } from "@/shared/ui";

import { emailError, passwordError } from "../lib/validation";
import { useAuth } from "../model/use-auth";

/** «С возвращением»: почта, пароль и ссылка на восстановление. */
export function LoginForm() {
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
    auth.mutate({ mode: "login", credentials: { email: email.trim(), password } });
  }

  return (
    <form onSubmit={submit} noValidate className="flex flex-col gap-7">
      <div className="flex flex-col gap-4">
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
            autoComplete="current-password"
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

        <Link
          href="/reset-password"
          className="self-end text-caption font-bold text-primary"
        >
          Забыли пароль?
        </Link>
      </div>

      <Button type="submit" disabled={auth.isPending}>
        {auth.isPending ? "Входим…" : "Войти"}
      </Button>
    </form>
  );
}
