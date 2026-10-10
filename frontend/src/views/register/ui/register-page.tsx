"use client";

import { RegisterForm, GuestGuard } from "@/features/auth";
import { GoalPicker, useOnboardingDraft } from "@/features/onboarding";
import { BetaBadge, Logo, Screen } from "@/shared/ui";

/** Шаг 2 из 2: аккаунт и ежедневная цель. */
export function RegisterPage() {
  const draft = useOnboardingDraft();

  return (
    <GuestGuard>
      <Screen className="px-gutter pt-6 pb-7">
        <header className="flex h-9 items-center justify-between">
          <Logo className="font-normal" />
          <BetaBadge />
        </header>

        <main className="mt-15.75 flex flex-1 flex-col gap-6.5">
          <div className="flex flex-col gap-2">
            <h1 className="text-display text-heading">
              Приобрети новую привычку.
              <br />
              Выучи казахский язык.
            </h1>
            <p className="text-body text-muted">
              Сохраняй ежедневный прогресс и увеличивай свои знания.
            </p>
          </div>

          <RegisterForm
            prefs={{
              dailyGoal: draft.dailyGoal,
              level: draft.level,
              motivation: draft.motivation,
            }}
          >
            <GoalPicker value={draft.dailyGoal} onChange={draft.setDailyGoal} />
          </RegisterForm>
        </main>
      </Screen>
    </GuestGuard>
  );
}
