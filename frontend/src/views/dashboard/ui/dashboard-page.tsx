"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { categoryQueries, type Category } from "@/entities/category";
import {
  levelAt,
  levelTitle,
  useLearnerStore,
  useLives,
} from "@/entities/learner";
import { Button, Eyebrow, StateNote } from "@/shared/ui";
import { LessonPreview } from "@/widgets/lesson-preview";
import { NoLivesDialog } from "@/widgets/no-lives";
import { LESSONS_PER_MODULE, SkillTree } from "@/widgets/skill-tree";

import { PathGlows } from "./path-glows";

/** Путь обучения: категории из API выстроены цепочкой уроков. */
export function DashboardPage() {
  const learner = useLearnerStore((state) => state.learner);
  const { lives } = useLives();
  const categories = useQuery(categoryQueries.all());
  const [preview, setPreview] = useState<{
    category: Category;
    index: number;
  } | null>(null);
  const [noLivesOpen, setNoLivesOpen] = useState(false);

  const completedIds = learner?.completedCategoryIds ?? [];
  const startLevel = learner?.level ?? "A1";
  const currentIndex = Math.max(
    0,
    categories.data?.findIndex(
      (category) => !completedIds.includes(category.id),
    ) ?? 0,
  );
  const currentLevel = levelAt(
    startLevel,
    Math.floor(currentIndex / LESSONS_PER_MODULE),
  );

  function select(category: Category, index: number) {
    if (lives === 0) setNoLivesOpen(true);
    else setPreview({ category, index });
  }

  return (
    <div className="relative flex flex-1 flex-col bg-background">
      <PathGlows />

      <div className="relative flex flex-1 flex-col px-gutter pt-6 pb-8">
        <div className="flex flex-col gap-1.25">
          <Eyebrow>{levelTitle(currentLevel)}</Eyebrow>
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
              completedIds={completedIds}
              startLevel={startLevel}
              onSelect={select}
            />
          )}
        </div>
      </div>

      <LessonPreview lesson={preview} onClose={() => setPreview(null)} />
      <NoLivesDialog open={noLivesOpen} onClose={() => setNoLivesOpen(false)} />
    </div>
  );
}
