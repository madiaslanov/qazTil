import { Fragment } from "react";

import type { PathLesson, PathUnit } from "@/entities/path";
import { cn } from "@/shared/lib/cn";
import { SkillNode } from "./skill-node";
import { UnitBanner, type UnitBannerState } from "./unit-banner";

/** Сдвиг узлов от центра, чтобы путь шёл змейкой, как в макете. */
const offsets = [
  "-translate-x-9.5",
  "translate-x-3.75",
  "-translate-x-13.5",
  "translate-x-11",
  "translate-x-0",
] as const;

function bannerState(unit: PathUnit): UnitBannerState {
  if (unit.lessons.every((lesson) => lesson.state === "completed")) return "completed";
  if (unit.lessons.every((lesson) => lesson.state === "locked")) return "locked";
  return "open";
}

/** Куда ведёт стрелка юнита: текущий урок, а у пройденного — первый. */
function entryLesson(unit: PathUnit): PathLesson | undefined {
  return (
    unit.lessons.find((lesson) => lesson.state === "current") ?? unit.lessons[0]
  );
}

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
          <li className={cn("w-full", unitIndex > 0 && "my-4")}>
            <UnitBanner
              number={unitIndex + 1}
              title={unit.title_ru}
              description={unit.description}
              state={bannerState(unit)}
              onStart={() => {
                const lesson = entryLesson(unit);
                if (lesson) onSelect(lesson, unit);
              }}
            />
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
