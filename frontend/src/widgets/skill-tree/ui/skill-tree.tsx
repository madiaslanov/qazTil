import { Fragment } from "react";

import type { Category } from "@/entities/category";
import { levelAt, type CourseLevel } from "@/entities/learner";
import { cn } from "@/shared/lib/cn";

import { ModuleTransition } from "./module-transition";
import { SkillNode, type SkillState } from "./skill-node";

/** Уроков в одном модуле: после них открывается следующий уровень. */
export const LESSONS_PER_MODULE = 5;

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

function chunk<T>(items: T[], size: number): T[][] {
  return Array.from({ length: Math.ceil(items.length / size) }, (_, index) =>
    items.slice(index * size, (index + 1) * size),
  );
}

/** Путь обучения: категории из API как цепочка уроков, разбитая на модули. */
export function SkillTree({
  categories,
  completedIds,
  startLevel,
  onSelect,
}: {
  categories: Category[];
  completedIds: number[];
  startLevel: CourseLevel;
  onSelect: (category: Category, index: number) => void;
}) {
  const current =
    categories.find((category) => !completedIds.includes(category.id)) ?? null;
  const modules = chunk(categories, LESSONS_PER_MODULE);

  return (
    <ol className="flex flex-col items-center gap-3">
      {modules.map((lessons, moduleIndex) => {
        const firstIndex = moduleIndex * LESSONS_PER_MODULE;
        const previous = modules[moduleIndex - 1];
        return (
          <Fragment key={lessons[0].id}>
            {previous && (
              <li className="my-4 w-full">
                <ModuleTransition
                  from={levelAt(startLevel, moduleIndex - 1)}
                  to={levelAt(startLevel, moduleIndex)}
                  unlocked={previous.every((lesson) =>
                    completedIds.includes(lesson.id),
                  )}
                  onStart={() => onSelect(lessons[0], firstIndex)}
                />
              </li>
            )}
            {lessons.map((category, offset) => {
              const index = firstIndex + offset;
              return (
                <li
                  key={category.id}
                  className={cn(offsets[index % offsets.length])}
                >
                  <SkillNode
                    label={category.name_ru}
                    state={stateOf(category, completedIds, current?.id ?? null)}
                    onSelect={() => onSelect(category, index)}
                  />
                </li>
              );
            })}
          </Fragment>
        );
      })}
    </ol>
  );
}
