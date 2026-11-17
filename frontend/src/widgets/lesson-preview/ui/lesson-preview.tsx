"use client";

import Link from "next/link";
import { BookOpen } from "lucide-react";

import type { PathLesson, PathUnit } from "@/entities/path";
import { BottomSheet, Button, Eyebrow } from "@/shared/ui";

/** Шторка урока: название, о чём он, и кнопка старта. */
export function LessonPreview({
  selected,
  onClose,
}: {
  selected: { lesson: PathLesson; unit: PathUnit } | null;
  onClose: () => void;
}) {
  return (
    <BottomSheet
      open={selected !== null}
      onOpenChange={(open) => !open && onClose()}
      title={selected?.lesson.title_ru ?? "Урок"}
    >
      {selected && <Content lesson={selected.lesson} unit={selected.unit} />}
    </BottomSheet>
  );
}

function Content({ lesson, unit }: { lesson: PathLesson; unit: PathUnit }) {
  const completed = lesson.state === "completed";

  return (
    <>
      <div className="flex items-start justify-between gap-4">
        <div className="flex flex-col gap-1.25">
          <Eyebrow className="text-micro">
            {unit.title_ru} · урок {lesson.position}
          </Eyebrow>
          <p className="text-h1 text-heading">{lesson.title_ru}</p>
          <p lang="kk" className="text-caption text-primary">
            {lesson.title_kk}
          </p>
        </div>
        <span className="flex size-13 shrink-0 items-center justify-center rounded-card bg-highlight-soft">
          <BookOpen className="size-6.25 text-foreground" />
        </span>
      </div>

      <p className="text-body text-muted">{lesson.description}</p>

      <Button asChild>
        <Link href={`/lesson/${lesson.id}`}>
          {completed ? "Пройти ещё раз" : "Начать"}
        </Link>
      </Button>
    </>
  );
}
