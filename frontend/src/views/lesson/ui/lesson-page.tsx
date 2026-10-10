"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { useLearnerStore, useLives } from "@/entities/learner";
import { QuizSession } from "@/features/quiz-session";
import type { Quiz } from "@/entities/quiz";
import { Screen } from "@/shared/ui";

import { LessonResult } from "./lesson-result";

/** Размер урока привязан к дневной цели: минута ≈ один вопрос. */
function questionsFor(goalMinutes: number | undefined) {
  return Math.max(4, Math.min(10, goalMinutes ?? 5));
}

export function LessonPage({ categoryId }: { categoryId: number }) {
  const router = useRouter();
  const dailyGoal = useLearnerStore((state) => state.learner?.dailyGoal);
  const { lives } = useLives();
  const [finished, setFinished] = useState<Quiz | null>(null);
  const [started] = useState(() => lives > 0);
  const onFinished = useCallback((quiz: Quiz) => setFinished(quiz), []);

  // Без жизней урок не начинаем — на пути покажется «Жизни закончились».
  useEffect(() => {
    if (!started) router.replace("/learn");
  }, [started, router]);

  if (!started) {
    return <Screen />;
  }

  if (finished) {
    return <LessonResult quiz={finished} />;
  }

  return (
    <Screen className="h-dvh bg-surface-muted">
      <QuizSession
        categoryId={categoryId}
        size={questionsFor(dailyGoal)}
        onFinished={onFinished}
      />
    </Screen>
  );
}
