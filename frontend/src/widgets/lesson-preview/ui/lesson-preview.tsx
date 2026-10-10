"use client";

import Link from "next/link";
import { BookOpen } from "lucide-react";
import { useQuery } from "@tanstack/react-query";

import type { Category } from "@/entities/category";
import { progressQueries } from "@/entities/progress";
import { wordQueries } from "@/entities/word";
import { BottomSheet, Button, Eyebrow, Progress } from "@/shared/ui";

/** Сколько новых слов показываем в превью. */
const PREVIEW_WORDS = 3;

/** Шторка урока: новые слова, прогресс темы и кнопка старта. */
export function LessonPreview({
  lesson,
  onClose,
}: {
  lesson: { category: Category; index: number } | null;
  onClose: () => void;
}) {
  return (
    <BottomSheet
      open={lesson !== null}
      onOpenChange={(open) => !open && onClose()}
      title={lesson?.category.name_ru ?? "Урок"}
    >
      {lesson && <Content category={lesson.category} index={lesson.index} />}
    </BottomSheet>
  );
}

function Content({ category, index }: { category: Category; index: number }) {
  const words = useQuery(wordQueries.list(category.id));
  const progress = useQuery(progressQueries.all());

  const topic = progress.data?.find((item) => item.category_id === category.id);
  const correct = topic?.correct_answers ?? 0;
  const total = topic?.total_answers ?? 0;
  const preview = words.data?.slice(0, PREVIEW_WORDS) ?? [];

  return (
    <>
      <div className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-1.25">
          <Eyebrow className="text-micro">
            Урок {index + 1}
            {words.data && ` · ${words.data.length} слов`}
          </Eyebrow>
          <p className="text-h1 text-heading">{category.name_ru}</p>
        </div>
        <span className="flex size-13 shrink-0 items-center justify-center rounded-card bg-highlight-soft">
          <BookOpen className="size-6.25 text-foreground" />
        </span>
      </div>

      <div className="flex flex-col gap-2.25">
        <p className="text-eyebrow font-extrabold text-muted uppercase">
          Новые слова · {words.data?.length ?? "…"}
        </p>
        <ul className="flex gap-2">
          {words.isPending &&
            Array.from({ length: PREVIEW_WORDS }, (_, key) => (
              <li
                key={key}
                className="h-13.5 flex-1 animate-pulse rounded-chip bg-background"
              />
            ))}
          {preview.map((word) => (
            <li
              key={word.id}
              className="flex min-w-0 flex-1 flex-col gap-0.75 rounded-chip bg-background p-2.5"
            >
              <span className="truncate text-caption text-primary">
                {word.kazakh}
              </span>
              <span className="truncate text-nano text-muted">
                {word.russian}
              </span>
            </li>
          ))}
        </ul>
      </div>

      <div className="flex flex-col gap-1.75">
        <div className="flex justify-between text-eyebrow">
          <span className="font-bold text-foreground">Прогресс темы</span>
          <span className="font-extrabold text-primary">
            {correct} из {total}
          </span>
        </div>
        <Progress
          tone="highlight"
          value={total ? Math.round((correct / total) * 100) : 0}
          aria-label="Верных ответов по теме"
        />
      </div>

      <Button asChild>
        <Link href={`/lesson/${category.id}`}>Начать</Link>
      </Button>
    </>
  );
}
