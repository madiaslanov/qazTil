"use client";

import Image from "next/image";

import { useQuery } from "@tanstack/react-query";

import { categoryQueries } from "@/entities/category";
import { useLearnerStore } from "@/entities/learner";
import { Button, StateNote } from "@/shared/ui";
import { SkillTree } from "@/widgets/skill-tree";

/** Путь обучения: категории из API выстроены цепочкой уроков. */
export function DashboardPage() {
  const completedIds = useLearnerStore(
    (state) => state.learner?.completedCategoryIds,
  );
  const categories = useQuery(categoryQueries.all());

  return (
    <div className="relative flex flex-1 flex-col">
      <Image
        src="/bg/dashboard.png"
        alt=""
        fill
        sizes="430px"
        className="pointer-events-none object-cover opacity-49 blur-[2px]"
      />

      <div className="relative flex flex-1 flex-col px-6 pt-6 pb-8">
        <p className="text-eyebrow font-extrabold uppercase text-muted">
          A1 · начальный
        </p>
        <h1 className="mt-1 text-h3">Твой путь обучения</h1>

        <div className="mt-5 flex-1">
          {categories.isPending && <StateNote text="Загружаем категории…" />}
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
              completedIds={completedIds ?? []}
            />
          )}
        </div>
      </div>
    </div>
  );
}
