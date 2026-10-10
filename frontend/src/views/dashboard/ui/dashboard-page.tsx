"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { categoryQueries, type Category } from "@/entities/category";
import { useLearnerStore, type Level } from "@/entities/learner";
import { Button, Eyebrow, StateNote } from "@/shared/ui";
import { LessonPreview } from "@/widgets/lesson-preview";
import { SkillTree } from "@/widgets/skill-tree";

import { PathGlows } from "./path-glows";

const levelTitles: Record<Level, string> = {
  A0: "A0 · с нуля",
  A1: "A1 · начальный",
  A2: "A2 · элементарный",
};

/** Путь обучения: категории из API выстроены цепочкой уроков. */
export function DashboardPage() {
  const learner = useLearnerStore((state) => state.learner);
  const categories = useQuery(categoryQueries.all());
  const [preview, setPreview] = useState<{
    category: Category;
    index: number;
  } | null>(null);

  return (
    <div className="relative flex flex-1 flex-col bg-background">
      <PathGlows />

      <div className="relative flex flex-1 flex-col px-gutter pt-6 pb-8">
        <div className="flex flex-col gap-1.25">
          <Eyebrow>{levelTitles[learner?.level ?? "A1"]}</Eyebrow>
          <h1 className="text-h3 text-foreground">Твой путь обучения</h1>
        </div>

        <div className="mt-10 flex flex-1 flex-col">
          {categories.isPending && <StateNote text="Загружаем уроки…" />}
          {categories.isError && (
            <StateNote
              text={categories.error.message}
              action={
                <Button
                  size="md"
                  variant="secondary"
                  className="w-auto"
                  onClick={() => categories.refetch()}
                >
                  Повторить
                </Button>
              }
            />
          )}
          {categories.data && (
            <SkillTree
              categories={categories.data}
              completedIds={learner?.completedCategoryIds ?? []}
              onSelect={(category, index) => setPreview({ category, index })}
            />
          )}
        </div>
      </div>

      <LessonPreview lesson={preview} onClose={() => setPreview(null)} />
    </div>
  );
}
