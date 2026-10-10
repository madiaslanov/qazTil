import type { Category } from "@/entities/category";
import { cn } from "@/shared/lib/cn";

import { SkillNode, type SkillState } from "./skill-node";

/** Сдвиг узлов от центра, чтобы путь шёл змейкой, как в макете. */
const offsets = [
  "-translate-x-9.5",
  "translate-x-3.75",
  "-translate-x-13.5",
  "translate-x-11",
  "translate-x-0",
] as const;

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
  onSelect,
}: {
  categories: Category[];
  completedIds: number[];
  onSelect: (category: Category, index: number) => void;
}) {
  const current =
    categories.find((category) => !completedIds.includes(category.id)) ?? null;

  return (
    <ol className="flex flex-col items-center gap-3">
      {categories.map((category, index) => (
        <li key={category.id} className={cn(offsets[index % offsets.length])}>
          <SkillNode
            label={category.name_ru}
            state={stateOf(category, completedIds, current?.id ?? null)}
            onSelect={() => onSelect(category, index)}
          />
        </li>
      ))}
    </ol>
  );
}
