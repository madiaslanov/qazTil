import { Check } from "lucide-react";

import type { Level } from "@/entities/learner";
import { cn } from "@/shared/lib/cn";
import { OptionCard } from "@/shared/ui";

const levels: { value: Level; title: string; note: string }[] = [
  { value: "A0", title: "Начинаю с нуля", note: "Нулевой уровень" },
  { value: "A1", title: "Знаю базовые фразы", note: "Начальный" },
  { value: "A2", title: "Могу поддержать диалог", note: "Элементарный" },
];

/** Стартовый уровень: три карточки с бейджем A0–A2. */
export function LevelPicker({
  value,
  onChange,
}: {
  value: Level;
  onChange: (level: Level) => void;
}) {
  return (
    <fieldset className="flex flex-col gap-2.5">
      <legend className="mb-2.5 text-caption font-extrabold text-primary">
        Твой уровень
      </legend>
      {levels.map((level) => {
        const selected = level.value === value;
        return (
          <OptionCard
            key={level.value}
            selected={selected}
            onClick={() => onChange(level.value)}
            className="flex h-20.5 items-center gap-3.5 px-4 text-left"
          >
            <span
              className={cn(
                "flex size-10.5 shrink-0 items-center justify-center rounded-chip text-body-sm font-extrabold",
                selected
                  ? "bg-primary text-primary-foreground"
                  : "bg-surface-muted text-primary",
              )}
            >
              {level.value}
            </span>
            <span className="flex flex-1 flex-col gap-0.75 font-bold">
              <span className="text-body">{level.title}</span>
              <span
                className={cn(
                  "text-micro uppercase",
                  selected ? "text-success" : "text-muted",
                )}
              >
                {selected ? "Выбрано" : level.note}
              </span>
            </span>
            {selected && <Check className="size-5 shrink-0 text-primary" />}
          </OptionCard>
        );
      })}
    </fieldset>
  );
}
