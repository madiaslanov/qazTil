"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { ChevronDown, CircleCheck, Lock, Play } from "lucide-react";

import {
  pathQueries,
  type LessonState,
  type Path,
  type PathUnit,
} from "@/entities/path";
import {
  progressQueries,
  type CourseProgress,
  type LessonProgress,
} from "@/entities/progress";
import { cn } from "@/shared/lib/cn";
import { Button, Card, Progress, StateNote } from "@/shared/ui";

const stateIcons = {
  completed: CircleCheck,
  current: Play,
  locked: Lock,
} as const;

const stateLabels: Record<LessonState, string> = {
  completed: "Пройден",
  current: "Идёт сейчас",
  locked: "Закрыт",
};

function percentOf(done: number, total: number) {
  return total === 0 ? 0 : Math.round((done / total) * 100);
}

/** Раскрыт по умолчанию модуль с текущим уроком, а если путь пройден — последний. */
function defaultUnitId(path: Path): number | null {
  const current = path.units.find((unit) =>
    unit.lessons.some((lesson) => lesson.id === path.current_lesson_id),
  );
  return current?.id ?? path.units.at(-1)?.id ?? null;
}

/** Прогресс по модулям. Модуль раскрывается в уроки с результатом по каждому. */
export function ProgressPage() {
  const path = useQuery(pathQueries.all());
  const progress = useQuery(progressQueries.all());
  // Что ученик открыл или закрыл сам; остальное решает defaultUnitId.
  const [toggled, setToggled] = useState<Record<number, boolean>>({});

  const failed = path.isError ? path : progress.isError ? progress : null;
  const ready = path.data && progress.data;
  const openByDefault = path.data ? defaultUnitId(path.data) : null;

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
      <div>
        <p className="text-eyebrow font-extrabold uppercase text-muted">
          Статистика
        </p>
        <h1 className="mt-1 text-h3">Твой прогресс</h1>
      </div>

      {!failed && !ready && <StateNote text="Считаем попытки…" />}
      {failed && (
        <StateNote
          text={failed.error?.message ?? "не удалось загрузить"}
          action={
            <Button
              size="md"
              variant="secondary"
              className="w-auto"
              onClick={() => {
                void path.refetch();
                void progress.refetch();
              }}
            >
              Повторить
            </Button>
          }
        />
      )}

      {ready && path.data.units.length === 0 && (
        <StateNote text="Пока ни одного урока — начни первый." />
      )}

      {ready &&
        path.data.units.map((unit, index) => (
          <UnitCard
            key={unit.id}
            unit={unit}
            number={index + 1}
            progress={progress.data}
            open={toggled[unit.id] ?? unit.id === openByDefault}
            onToggle={(open) =>
              setToggled((previous) => ({ ...previous, [unit.id]: open }))
            }
          />
        ))}
    </div>
  );
}

function UnitCard({
  unit,
  number,
  progress,
  open,
  onToggle,
}: {
  unit: PathUnit;
  number: number;
  progress: CourseProgress;
  open: boolean;
  onToggle: (open: boolean) => void;
}) {
  const stats = progress.units.find((item) => item.unit_id === unit.id);
  const completed = stats?.completed_lessons ?? 0;
  const total = stats?.total_lessons ?? unit.lessons.length;
  const percent = percentOf(completed, total);
  const byLesson = new Map<number, LessonProgress>(
    progress.lessons.map((lesson) => [lesson.lesson_id, lesson]),
  );
  const panel = `unit-${unit.id}`;

  return (
    <Card className="flex flex-col p-0">
      <button
        type="button"
        onClick={() => onToggle(!open)}
        aria-expanded={open}
        aria-controls={panel}
        className="flex flex-col gap-3 rounded-[inherit] p-5 text-left outline-none"
      >
        <div className="flex items-center justify-between gap-3">
          <div className="min-w-0">
            <p className="text-micro font-extrabold uppercase text-muted">
              Модуль {number}
            </p>
            <p className="truncate text-stat">{unit.title_ru}</p>
          </div>
          <div className="flex shrink-0 items-center gap-2">
            <b className="text-stat font-extrabold">{percent}%</b>
            <ChevronDown
              className={cn(
                "size-5 text-muted transition-transform",
                open && "rotate-180",
              )}
            />
          </div>
        </div>
        <Progress value={percent} tone="highlight" />
        <p className="text-caption text-muted">
          Уроков {completed} из {total}
        </p>
      </button>

      {open && (
        <ul id={panel} className="flex flex-col gap-1 border-t border-line px-5 py-3">
          {unit.lessons.map((lesson) => {
            const result = byLesson.get(lesson.id);
            const Icon = stateIcons[lesson.state];
            const answered = result?.total
              ? `${result.correct} из ${result.total}`
              : null;
            return (
              <li key={lesson.id} className="flex items-center gap-3 py-2">
                <Icon
                  className={cn(
                    "size-5 shrink-0",
                    lesson.state === "completed" && "text-success",
                    lesson.state === "current" && "text-primary",
                    lesson.state === "locked" && "text-muted",
                  )}
                />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-body text-foreground">
                    {lesson.title_ru}
                  </p>
                  <p className="text-caption text-muted">
                    {stateLabels[lesson.state]}
                  </p>
                </div>
                {answered && lesson.state !== "locked" && (
                  <span className="shrink-0 text-caption font-bold text-muted tabular-nums">
                    {answered}
                  </span>
                )}
              </li>
            );
          })}
        </ul>
      )}
    </Card>
  );
}
