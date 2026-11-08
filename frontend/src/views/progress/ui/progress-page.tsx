"use client";

import { useQuery } from "@tanstack/react-query";

import { progressQueries } from "@/entities/progress";
import { Button, Card, Progress, StateNote } from "@/shared/ui";

/** Прогресс по урокам и модулям: его считает сервер по попыткам. */
export function ProgressPage() {
  const progress = useQuery(progressQueries.all());
  const lessons = progress.data?.lessons ?? [];
  const units = progress.data?.units ?? [];

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
      <div>
        <p className="text-eyebrow font-extrabold uppercase text-muted">
          Статистика
        </p>
        <h1 className="mt-1 text-h3">Твой прогресс</h1>
      </div>

      {progress.isPending && <StateNote text="Считаем попытки…" />}
      {progress.isError && (
        <StateNote
          text={progress.error.message}
          action={
            <Button
              size="md"
              variant="secondary"
              className="w-auto"
              onClick={() => progress.refetch()}
            >
              Повторить
            </Button>
          }
        />
      )}

      {progress.data && lessons.length === 0 && units.length === 0 && (
        <StateNote text="Пока ни одного урока — начни первый." />
      )}

      {units.map((unit) => {
        const percent =
          unit.total_lessons === 0
            ? 0
            : Math.round((unit.completed_lessons / unit.total_lessons) * 100);
        return (
          <Card key={unit.unit_id} className="flex flex-col gap-3 p-5">
            <div className="flex items-baseline justify-between gap-3">
              <p className="text-stat">Модуль {unit.unit_id}</p>
              <b className="text-stat font-extrabold">{percent}%</b>
            </div>
            <Progress value={percent} tone="highlight" />
            <p className="text-caption text-muted">
              Уроков {unit.completed_lessons} из {unit.total_lessons}
            </p>
          </Card>
        );
      })}

      {lessons.map((lesson) => {
        const percent =
          lesson.total === 0
            ? 0
            : Math.round((lesson.correct / lesson.total) * 100);
        return (
          <Card key={lesson.lesson_id} className="flex flex-col gap-3 p-5">
            <div className="flex items-baseline justify-between gap-3">
              <div>
                <p className="text-stat">Урок {lesson.lesson_id}</p>
                <p className="text-caption text-muted">
                  {lesson.completed ? "Пройден" : "Ещё не пройден"}
                </p>
              </div>
              <b className="text-stat font-extrabold">{percent}%</b>
            </div>
            <Progress value={percent} tone="highlight" />
            <p className="text-caption text-muted">
              Прогресс темы {lesson.correct} из {lesson.total}
            </p>
          </Card>
        );
      })}
    </div>
  );
}
