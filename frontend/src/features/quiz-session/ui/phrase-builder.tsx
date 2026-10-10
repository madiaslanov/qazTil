"use client";

import type { AssembleExercise } from "@/entities/exercise";

import { WordChip } from "./word-chip";

/** «Собери фразу»: плитки из банка переносятся в ответ и обратно. */
export function PhraseBuilder({
  exercise,
  value,
  disabled,
  onChange,
}: {
  exercise: AssembleExercise;
  value: string[];
  disabled: boolean;
  /** Обновление через функцию: несколько нажатий подряд не теряются. */
  onChange: (update: (tileIds: string[]) => string[]) => void;
}) {
  const byId = new Map(exercise.tiles.map((tile) => [tile.id, tile]));

  return (
    <>
      <div className="relative flex min-h-47.5 flex-col rounded-card bg-background px-3.5 py-4.5">
        <div aria-hidden className="pointer-events-none absolute inset-x-0 top-21.5 flex flex-col gap-5">
          <span className="h-0.5 rounded-full bg-primary" />
          <span className="h-0.5 rounded-full bg-line" />
          <span className="h-0.5 rounded-full bg-line" />
        </div>
        <div className="relative flex flex-wrap gap-2" aria-label="Твой ответ">
          {value.map((id) => (
            <WordChip
              key={id}
              text={byId.get(id)?.text ?? ""}
              placed
              disabled={disabled}
              onClick={() => onChange((ids) => ids.filter((item) => item !== id))}
            />
          ))}
        </div>
      </div>

      <div className="flex flex-col gap-3.5">
        <div className="flex justify-between text-micro text-muted">
          <span className="font-extrabold">СЛОВА</span>
          <span>Нажми, чтобы добавить</span>
        </div>
        <div className="flex flex-wrap gap-2.25">
          {exercise.tiles.map((tile) => (
            <WordChip
              key={tile.id}
              text={tile.text}
              used={value.includes(tile.id)}
              disabled={disabled}
              onClick={() => onChange((ids) => [...ids, tile.id])}
            />
          ))}
        </div>
      </div>
    </>
  );
}
