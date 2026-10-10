"use client";

import Link from "next/link";

import { GuestGuard, LoginForm } from "@/features/auth";
import { BackLink, Screen } from "@/shared/ui";

/** «С возвращением» для тех, у кого уже есть аккаунт. */
export function LoginPage() {
  return (
    <GuestGuard>
      <Screen className="px-gutter pt-6 pb-7">
        <header className="flex h-9 items-center justify-between">
          <BackLink href="/" />
          <Link href="/setup" className="text-caption font-bold text-accent">
            Регистрация
          </Link>
        </header>

        <main className="mt-18 flex flex-col gap-7">
          <div className="flex flex-col gap-2.25">
            <h1 className="text-display font-semibold text-heading">
              С возвращением
            </h1>
            <p className="text-body text-muted">
              Продолжи свой путь с того места, где остановился.
            </p>
          </div>
          <LoginForm />
        </main>
      </Screen>
    </GuestGuard>
  );
}
