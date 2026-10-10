"use client";

import { useState } from "react";
import { CircleAlert, CircleCheck, MailCheck } from "lucide-react";

import { Button, FormField, Input } from "@/shared/ui";

import { emailError } from "../lib/validation";

/**
 * Запрос ссылки на новый пароль. В Go API восстановления пока нет,
 * поэтому после отправки честно говорим, что функция в работе.
 */
export function ResetPasswordForm() {
  const [email, setEmail] = useState("");
  const [touched, setTouched] = useState(false);
  const [sent, setSent] = useState(false);

  const problem = emailError(email);
  const status = !touched ? "default" : problem ? "error" : "success";

  function submit(event: React.FormEvent) {
    event.preventDefault();
    setTouched(true);
    if (!problem) setSent(true);
  }

  return (
    <form onSubmit={submit} noValidate className="flex flex-col gap-7">
      <FormField
        id="email"
        label="Почта"
        hint={touched && (problem ?? "Адрес выглядит верно")}
        hintTone={problem ? "error" : "success"}
      >
        <Input
          id="email"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          value={email}
          status={status}
          icon={
            status === "error" ? (
              <CircleAlert />
            ) : status === "success" ? (
              <CircleCheck />
            ) : null
          }
          aria-describedby={touched ? "email-hint" : undefined}
          onChange={(event) => {
            setEmail(event.target.value);
            setSent(false);
          }}
          onBlur={() => setTouched(true)}
        />
      </FormField>

      <Button type="submit">Отправить ссылку</Button>

      {sent && (
        <div
          role="status"
          className="flex items-start gap-3 rounded-control bg-success-soft p-4"
        >
          <MailCheck className="size-5.5 shrink-0 text-success" />
          <div className="flex flex-col gap-1">
            <p className="text-body-sm font-extrabold text-primary">
              Восстановление скоро заработает
            </p>
            <p className="text-eyebrow text-muted">
              Сервер пока не отправляет письма — эту часть доделаем отдельной
              задачей.
            </p>
          </div>
        </div>
      )}
    </form>
  );
}
