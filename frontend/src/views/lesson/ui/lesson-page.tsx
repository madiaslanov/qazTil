"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";

import { profileQueries, useLives } from "@/entities/learner";
import { QuizSession, type LessonScore } from "@/features/quiz-session";
import { Screen } from "@/shared/ui";

import { LessonResult } from "./lesson-result";

/** Размер урока привязан к дневной цели: минута ≈ один вопрос. */
function questionsFor(goalMinutes: number | undefined) {
  return Math.max(4, Math.min(10, goalMinutes ?? 5));
}

export function LessonPage({ categoryId }: { categoryId: number }) {
  const router = useRouter();
  const profile = useQuery(profileQueries.current());
  const { lives, pending } = useLives();
  const [finished, setFinished] = useState<LessonScore | null>(null);
  const onFinished = useCallback(
    (score: LessonScore) => setFinished(score),
    [],
  );

  // Без жизней урок не начинаем — на пути покажется «Жизни закончились».
  useEffect(() => {
    if (!pending && lives === 0) router.replace("/learn");
  }, [lives, pending, router]);

  if (pending || profile.isPending || lives === 0) {
    return <Screen />;
  }

  if (finished) {
    return <LessonResult score={finished} />;
  }

  return (
    <Screen className="h-dvh bg-surface-muted">
      <QuizSession
        categoryId={categoryId}
        size={questionsFor(profile.data?.daily_goal)}
        onFinished={onFinished}
      />
    </Screen>
  );
}
