"use client";

import Link from "next/link";

import { GuestGuard } from "@/features/auth";
import {
  LevelPicker,
  MotivationPicker,
  useOnboardingDraft,
} from "@/features/onboarding";
import { BackLink, Button, Screen } from "@/shared/ui";

/** Шаг 1 из 2: уровень и причина учить язык. */
export function SetupPage() {
  const draft = useOnboardingDraft();

  return (
    <GuestGuard>
      <Screen className="px-gutter pt-6 pb-7">
        <header className="flex h-9 items-center justify-between">
          <BackLink href="/" />
          <p className="text-caption font-bold text-accent">Шаг 1 из 2</p>
        </header>

        <main className="mt-7 flex flex-col gap-5">
          <div className="flex flex-col gap-2">
            <h1 className="text-h1 font-semibold text-heading">
              Настроим твой маршрут
            </h1>
            <p className="text-body-sm text-muted">
              Выбери текущий уровень и главную причину учить казахский.
            </p>
          </div>
          <LevelPicker value={draft.level} onChange={draft.setLevel} />
          <MotivationPicker
            value={draft.motivation}
            onChange={draft.setMotivation}
          />
        </main>

        <Button asChild className="mt-auto">
          <Link href="/register">Продолжить</Link>
        </Button>
      </Screen>
    </GuestGuard>
  );
}
