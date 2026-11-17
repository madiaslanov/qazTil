"use client";

import { useState } from "react";

import type { MatchExercise } from "@/entities/exercise";
import { cn } from "@/shared/lib/cn";

type Side = "left" | "right";
type Pair = [string, string];

function plural(count: number) {
  if (count % 10 === 1 && count % 100 !== 11) return "пара собрана";
  if ([2, 3, 4].includes(count % 10) && ![12, 13, 14].includes(count % 100)) {
    return "пары собраны";
  }
  return "пар собрано";
}

/**
 * «Найди пару»: нажми слово и его перевод — они связываются номером.
 * Верна ли связь, сервер скажет только на весь набор сразу, поэтому доска
 * ничего не подсвечивает как правильное. Нажатие на связанную плитку
 * разрывает пару.
 */
export function MatchBoard({
  exercise,
  pairs,
  disabled,
  onChange,
}: {
  exercise: MatchExercise;
  pairs: Pair[];
  disabled: boolean;
  onChange: (pairs: Pair[]) => void;
}) {
  const [picked, setPicked] = useState<{ side: Side; id: string } | null>(null);

  const numberOf = (id: string) => {
    const at = pairs.findIndex((pair) => pair.includes(id));
    return at === -1 ? null : at + 1;
  };

  function pick(side: Side, id: string) {
    if (disabled) return;
    if (numberOf(id) !== null) {
      onChange(pairs.filter((pair) => !pair.includes(id)));
      setPicked(null);
      return;
    }
    if (!picked || picked.side === side) {
      setPicked({ side, id });
      return;
    }
    const left = side === "left" ? id : picked.id;
    const right = side === "right" ? id : picked.id;
    setPicked(null);
    onChange([...pairs, [left, right]]);
  }

  const rows = exercise.pairs_left.map((left, index) => [
    { side: "left" as const, item: left },
    { side: "right" as const, item: exercise.pairs_right[index] },
  ]);

  return (
    <div className="flex flex-col gap-3">
      <p className="flex items-center gap-1.5 text-micro font-bold text-success">
        <span aria-hidden className="size-1.5 rounded-full bg-success" />
        {pairs.length} {plural(pairs.length)}
      </p>
      <div className="grid grid-cols-2 gap-2.5">
        {rows.flat().map(({ side, item }) => {
          const number = numberOf(item.id);
          const active = picked?.id === item.id;
          return (
            <button
              key={item.id}
              type="button"
              lang={side === "left" ? "kk" : undefined}
              onClick={() => pick(side, item.id)}
              disabled={disabled}
              aria-pressed={active || number !== null}
              className={cn(
                "relative flex h-17 items-center justify-center rounded-control border px-3 text-center text-body text-primary transition-colors outline-none",
                "border-line bg-surface",
                active && "border-2 border-primary bg-highlight-soft",
                number !== null && "border-2 border-highlight bg-highlight-soft",
              )}
            >
              {number !== null && (
                <span className="absolute top-1.5 left-2 text-micro font-extrabold text-muted">
                  {number}
                </span>
              )}
              {item.text}
            </button>
          );
        })}
      </div>
    </div>
  );
}
