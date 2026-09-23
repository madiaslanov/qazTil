"use client";

import { useCallback, useState } from "react";
import Image from "next/image";

import { useLearnerStore } from "@/entities/learner";
import { QuizSession } from "@/features/quiz-session";
import type { Quiz } from "@/entities/quiz";
import { Screen } from "@/shared/ui";

import { LessonResult } from "./lesson-result";

/** Размер урока привязан к дневной цели: минута ≈ один вопрос. */
function questionsFor(goalMinutes: number | undefined) {
  return Math.max(4, Math.min(10, goalMinutes ?? 5));
}

export function LessonPage({ categoryId }: { categoryId: number }) {
  const dailyGoal = useLearnerStore((state) => state.learner?.dailyGoal);
  const [finished, setFinished] = useState<Quiz | null>(null);
  const onFinished = useCallback((quiz: Quiz) => setFinished(quiz), []);

  if (finished) {
    return <LessonResult quiz={finished} />;
  }

  return (
    <Screen>
      <Image
        src="/bg/lesson.png"
        alt=""
        fill
        sizes="430px"
        className="pointer-events-none object-cover mix-blend-color-burn"
      />
      <div className="relative flex min-h-0 flex-1 flex-col">
        <QuizSession
          categoryId={categoryId}
          size={questionsFor(dailyGoal)}
          onFinished={onFinished}
        />
      </div>
    </Screen>
  );
}
