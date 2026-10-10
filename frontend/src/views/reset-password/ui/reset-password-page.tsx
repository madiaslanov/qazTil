"use client";

import { KeyRound } from "lucide-react";

import { GuestGuard, ResetPasswordForm } from "@/features/auth";
import { BackLink, Screen } from "@/shared/ui";

/** «Вернем доступ»: ссылка на новый пароль по почте. */
export function ResetPasswordPage() {
  return (
    <GuestGuard>
      <Screen className="px-gutter pt-6 pb-7">
        <header className="flex h-9 items-center">
          <BackLink href="/login" />
        </header>

        <main className="mt-18 flex flex-col gap-7">
          <span className="flex size-17 items-center justify-center rounded-card bg-highlight-soft">
            <KeyRound className="size-7.5 text-primary" />
          </span>
          <div className="flex flex-col gap-2.25">
            <h1 className="text-display font-semibold text-heading">
              Вернем доступ
            </h1>
            <p className="text-body text-muted">
              Отправим ссылку для нового пароля на подтвержденную почту.
            </p>
          </div>
          <ResetPasswordForm />
        </main>
      </Screen>
    </GuestGuard>
  );
}
