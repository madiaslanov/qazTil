"use client";

import Image from "next/image";

import { categoryApi } from "@/entities/category";
import { useLearner } from "@/entities/learner";
import { useRequest } from "@/shared/lib/use-request";
import { Button, Screen, StateNote } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { SkillTree } from "@/widgets/skill-tree";
import { TopBar } from "@/widgets/top-bar";

/** Путь обучения: категории из API выстроены цепочкой уроков. */
export function DashboardPage() {
  const learner = useLearner();
  const categories = useRequest(() => categoryApi.list(), "categories");

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
            {categories.loading && <StateNote text="Загружаем категории…" />}
            {categories.error && (
              <StateNote
                text={categories.error}
                action={
                  <Button
                    size="md"
                    variant="quiet"
                    className="w-auto"
                    onClick={categories.retry}
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
              />
            )}
          </div>
        </div>
      </div>

      <BottomNav />
    </Screen>
  );
}
