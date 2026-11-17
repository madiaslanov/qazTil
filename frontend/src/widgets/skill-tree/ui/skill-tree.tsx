import { Fragment } from "react";

import type { PathLesson, PathUnit } from "@/entities/path";
import { cn } from "@/shared/lib/cn";
import { Eyebrow } from "@/shared/ui";

import { SkillNode } from "./skill-node";

/** Сдвиг узлов от центра, чтобы путь шёл змейкой, как в макете. */
const offsets = [
  "-translate-x-9.5",
  "translate-x-3.75",
  "-translate-x-13.5",
  "translate-x-11",
  "translate-x-0",
] as const;

/**
 * Путь обучения: юниты и уроки приходят из GET /path уже по порядку,
 * а состояние каждого урока считает сервер.
 */
export function SkillTree({
  units,
  onSelect,
}: {
  units: PathUnit[];
  onSelect: (lesson: PathLesson, unit: PathUnit) => void;
}) {
  let nodeIndex = 0;

  return (
    <ol className="flex flex-col items-center gap-3">
      {units.map((unit, unitIndex) => (
        <Fragment key={unit.id}>
          <li className={cn("w-full", unitIndex > 0 && "mt-4")}>
            <div className="flex flex-col gap-1 rounded-card border-2 border-surface bg-surface-muted px-5 py-4">
              <Eyebrow className="text-micro">Юнит {unitIndex + 1}</Eyebrow>
              <p className="text-title text-primary">{unit.title_ru}</p>
              <p className="text-caption text-muted">{unit.description}</p>
            </div>
          </li>
          {unit.lessons.map((lesson) => {
            const offset = offsets[nodeIndex++ % offsets.length];
            return (
              <li key={lesson.id} className={cn(offset)}>
                <SkillNode
                  label={lesson.title_ru}
                  state={lesson.state}
                  onSelect={() => onSelect(lesson, unit)}
                />
              </li>
            );
          })}
        </Fragment>
      ))}
    </ol>
  );
}
