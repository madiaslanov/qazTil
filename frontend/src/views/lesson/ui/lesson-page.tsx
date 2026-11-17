"use client";

import { useCallback, useState } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { lessonQueries } from "@/entities/lesson";
import { LessonSession, type LessonSummary } from "@/features/lesson-session";
import { ApiError } from "@/shared/api";
import { Button, Screen, StateNote } from "@/shared/ui";

import { LessonResult } from "./lesson-result";

function BackToPath({ text }: { text: string }) {
  return (
    <Screen>
      <StateNote
        text={text}
        action={
          <Button asChild size="md" variant="secondary" className="w-auto">
            <Link href="/learn">Вернуться к пути</Link>
          </Button>
        }
      />
    </Screen>
  );
}

export function LessonPage({ lessonId }: { lessonId: number }) {
  const lesson = useQuery(lessonQueries.detail(lessonId));
  const [finished, setFinished] = useState<LessonSummary | null>(null);
  const onFinished = useCallback(
    (summary: LessonSummary) => setFinished(summary),
    [],
  );

  if (lesson.isPending) return <Screen />;

  if (lesson.isError) {
    const status = lesson.error instanceof ApiError ? lesson.error.status : 0;
    if (status === 403) {
      return <BackToPath text="Этот урок ещё закрыт: пройди предыдущие." />;
    }
    if (status === 404) return <BackToPath text="Такого урока нет." />;
    return <BackToPath text={lesson.error.message} />;
  }

  // Без жизней урок не начинаем. Если жизни есть, но их 0 по данным урока —
  // сервер всё равно ответит 409 на первый ответ, но не пускаем зря.
  if (lesson.data.lives === 0) {
    return <BackToPath text="Жизни закончились. Вернись, когда жизнь восстановится." />;
  }

  if (finished) {
    return <LessonResult summary={finished} />;
  }

  return (
    <Screen className="h-dvh bg-surface-muted">
      <LessonSession lesson={lesson.data} onFinished={onFinished} />
    </Screen>
  );
}
