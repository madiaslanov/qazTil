"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { learnerStore, type DailyGoal } from "@/entities/learner";
import { Button, Input, Label } from "@/shared/ui";

import { GoalPicker } from "./goal-picker";

/**
 * Вход из макета. Авторизации в API нет, поэтому пароль никуда не уходит,
 * а почта и цель просто заводят локальный профиль ученика.
 */
export function OnboardingForm() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [goal, setGoal] = useState<DailyGoal>(10);

  function start(event: React.FormEvent) {
    event.preventDefault();
    learnerStore.create(email.trim() || "ученик", goal);
    router.push("/learn");
  }

  return (
    <form onSubmit={start} className="flex flex-1 flex-col gap-6.5">
      <div className="flex flex-col gap-3.5">
        <div className="flex flex-col gap-[7px]">
          <Label htmlFor="email">Почта</Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
        </div>
        <div className="flex flex-col gap-[7px]">
          <Label htmlFor="password">Пароль</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            placeholder="••••••••"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
        </div>
      </div>

      <GoalPicker value={goal} onChange={setGoal} />

      <div className="mt-auto pt-6">
        <Button type="submit" variant="primary">
          Начать обучение
        </Button>
      </div>
    </form>
  );
}
