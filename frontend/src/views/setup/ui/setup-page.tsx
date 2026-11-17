"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  profileApi,
  profileQueries,
  type DailyGoal,
  type Level,
  type Motivation,
} from "@/entities/learner";
import { useSessionStore } from "@/entities/session";
import {
  GoalPicker,
  LevelPicker,
  MotivationPicker,
  useOnboardingDraft,
} from "@/features/onboarding";
import { BackLink, Button, FieldHint, Screen } from "@/shared/ui";

/**
 * Настройка маршрута. Гость заполняет первый шаг и уходит регистрироваться,
 * вошедший ученик отвечает здесь же и сохраняет ответы на сервере: до
 * регистрации профиля ещё нет, после — он уже есть и его надо править.
 */
export function SetupPage() {
  const session = useSessionStore((state) => state.session);
  const hydrated = useSessionStore((state) => state.hydrated);

  if (!hydrated) {
    return <Screen />;
  }
  return session ? <SignedInSetup /> : <GuestSetup />;
}

/** Шаг 1 из 2: уровень и причина учить язык. Цель спросят при регистрации. */
function GuestSetup() {
  const draft = useOnboardingDraft();

  return (
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
  );
}

/**
 * Все три вопроса сразу: аккаунт уже есть, второго шага не будет.
 * Форма ждёт профиль и только потом монтируется, чтобы поля сразу завелись
 * сохранёнными ответами, а не значениями по умолчанию.
 */
function SignedInSetup() {
  const profile = useQuery(profileQueries.current());

  if (profile.isPending || !profile.data) {
    return <Screen />;
  }
  return (
    <SetupForm
      level={profile.data.level}
      motivation={profile.data.motivation}
      dailyGoal={profile.data.daily_goal}
    />
  );
}

function SetupForm(initial: {
  level: Level;
  motivation: Motivation;
  dailyGoal: DailyGoal;
}) {
  const router = useRouter();
  const queryClient = useQueryClient();

  const [level, setLevel] = useState<Level>(initial.level);
  const [motivation, setMotivation] = useState<Motivation>(initial.motivation);
  const [dailyGoal, setDailyGoal] = useState<DailyGoal>(initial.dailyGoal);

  const save = useMutation({
    mutationFn: () =>
      profileApi.update({ level, motivation, daily_goal: dailyGoal }),
    onSuccess: (data) => {
      queryClient.setQueryData(profileQueries.current().queryKey, data);
      router.replace("/learn");
    },
  });

  return (
    <Screen className="px-gutter pt-6 pb-7">
      <header className="flex h-9 items-center justify-between">
        <p className="text-caption font-bold text-accent">Настройка маршрута</p>
      </header>

      <main className="mt-7 flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <h1 className="text-h1 font-semibold text-heading">
            Настроим твой маршрут
          </h1>
          <p className="text-body-sm text-muted">
            Ответь на три вопроса — подберём уроки и дневную норму.
          </p>
        </div>
        <LevelPicker value={level} onChange={setLevel} />
        <MotivationPicker value={motivation} onChange={setMotivation} />
        <GoalPicker value={dailyGoal} onChange={setDailyGoal} />
        {save.isError ? (
          <FieldHint tone="error">
            Не получилось сохранить. Попробуй ещё раз.
          </FieldHint>
        ) : null}
      </main>

      <Button
        className="mt-auto"
        disabled={save.isPending}
        onClick={() => save.mutate()}
      >
        {save.isPending ? "Сохраняем…" : "Начать учиться"}
      </Button>
    </Screen>
  );
}
