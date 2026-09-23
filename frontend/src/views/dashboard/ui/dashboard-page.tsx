"use client";

import Image from "next/image";

import { useQuery } from "@tanstack/react-query";

import { categoryQueries } from "@/entities/category";
import { useLearnerStore } from "@/entities/learner";
import { Button, Screen, StateNote } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { SkillTree } from "@/widgets/skill-tree";
import { TopBar } from "@/widgets/top-bar";

/** Путь обучения: категории из API выстроены цепочкой уроков. */
export function DashboardPage() {
  const completedIds = useLearnerStore(
    (state) => state.learner?.completedCategoryIds,
  );
  const categories = useQuery(categoryQueries.all());

  return (
    <Screen>
      <TopBar />

      <div className="relative flex min-h-0 flex-1 flex-col overflow-y-auto">
        <Image
          src="/bg/dashboard.png"
          alt=""
          fill
          sizes="430px"
          className="pointer-events-none object-cover opacity-49 blur-[2px]"
        />

        <div className="relative flex flex-1 flex-col px-6 pt-6 pb-8">
          <p className="text-[12px] font-extrabold uppercase text-muted">
            A1 · начальный
          </p>
          <h1 className="mt-1 text-[26px]">Твой путь обучения</h1>

          <div className="mt-[20px] flex-1">
            {categories.isPending && <StateNote text="Загружаем категории…" />}
            {categories.isError && (
              <StateNote
                text={categories.error.message}
                action={
                  <Button
                    size="md"
                    variant="quiet"
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

      <BottomNav />
    </Screen>
  );
}
