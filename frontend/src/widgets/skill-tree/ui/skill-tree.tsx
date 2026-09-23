"use client";

import type { Category } from "@/entities/category";
import { cn } from "@/shared/lib/cn";

import { SkillNode, type SkillState } from "./skill-node";

/** Смещения узлов по горизонтали, чтобы путь шёл змейкой, как в макете. */
const offsets = [
  "justify-center",
  "justify-start pl-[72px]",
  "justify-end pr-[62px]",
  "justify-center",
] as const;

const connectors = ["rotate-0", "-rotate-25", "rotate-24", "rotate-0"] as const;

function stateOf(
  category: Category,
  completedIds: number[],
  currentId: number | null,
): SkillState {
  if (completedIds.includes(category.id)) return "completed";
  if (category.id === currentId) return "current";
  return "locked";
}

/** Путь обучения: категории из API как цепочка уроков. */
export function SkillTree({
  categories,
  completedIds,
}: {
  categories: Category[];
  completedIds: number[];
}) {
  const current =
    categories.find((category) => !completedIds.includes(category.id)) ?? null;

  return (
    <div className="flex flex-col items-stretch">
      {categories.map((category, index) => {
        const state = stateOf(category, completedIds, current?.id ?? null);
        return (
          <div key={category.id} className="flex flex-col items-center">
            {index > 0 && (
              <span
                className={cn(
                  "my-1 h-[24px] w-[4px] rounded-full",
                  connectors[index % connectors.length],
                  state === "locked" ? "bg-track" : "bg-ink",
                )}
              />
            )}
            <div
              className={cn("flex w-full", offsets[index % offsets.length])}
            >
              <SkillNode
                href={`/lesson/${category.id}`}
                label={category.name_ru}
                state={state}
              />
            </div>
          </div>
        );
      })}
    </div>
  );
}
