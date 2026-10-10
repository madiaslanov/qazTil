"use client";

import { useEffect, useRef, useState } from "react";

import { isLocalPair, type MatchExercise } from "@/entities/exercise";
import { cn } from "@/shared/lib/cn";

type Side = "left" | "right";

/** Сколько держим коралловую подсветку неверной пары. */
const MISS_MS = 600;

function plural(count: number) {
  if (count % 10 === 1 && count % 100 !== 11) return "пара найдена";
  if ([2, 3, 4].includes(count % 10) && ![12, 13, 14].includes(count % 100)) {
    return "пары найдены";
  }
  return "пар найдено";
}

/**
 * «Найди пару»: нажми слово и его перевод. Верная пара становится мятной,
 * неверная коротко мигает коралловым. Когда найдены все — onComplete.
 */
export function MatchBoard({
  exercise,
  disabled,
  onComplete,
}: {
  exercise: MatchExercise;
  disabled: boolean;
  onComplete: (pairs: [string, string][]) => void;
}) {
  const [picked, setPicked] = useState<{ side: Side; id: string } | null>(null);
  const [matched, setMatched] = useState<[string, string][]>([]);
  const [missed, setMissed] = useState<string[]>([]);
  const missTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => () => {
    if (missTimer.current) clearTimeout(missTimer.current);
  }, []);

  const isMatched = (id: string) => matched.some((pair) => pair.includes(id));

  function pick(side: Side, id: string) {
    if (disabled || isMatched(id)) return;
    if (!picked || picked.side === side) {
      setPicked({ side, id });
      return;
    }
    const left = side === "left" ? id : picked.id;
    const right = side === "right" ? id : picked.id;
    setPicked(null);

    if (isLocalPair(exercise, left, right)) {
      const next: [string, string][] = [...matched, [left, right]];
      setMatched(next);
      if (next.length === exercise.pairs_left.length) onComplete(next);
      return;
    }
    setMissed([left, right]);
    if (missTimer.current) clearTimeout(missTimer.current);
    missTimer.current = setTimeout(() => setMissed([]), MISS_MS);
  }

  const rows = exercise.pairs_left.map((left, index) => [
    { side: "left" as const, item: left },
    { side: "right" as const, item: exercise.pairs_right[index] },
  ]);

  return (
    <div className="flex flex-col gap-3">
      <p className="flex items-center gap-1.5 text-micro font-bold text-success">
        <span aria-hidden className="size-1.5 rounded-full bg-success" />
        {matched.length} {plural(matched.length)}
      </p>
      <div className="grid grid-cols-2 gap-2.5">
        {rows.flat().map(({ side, item }) => {
          const done = isMatched(item.id);
          const active = picked?.id === item.id;
          const miss = missed.includes(item.id);
          return (
            <button
              key={item.id}
              type="button"
              lang={side === "left" ? "kk" : undefined}
              onClick={() => pick(side, item.id)}
              disabled={disabled || done}
              aria-pressed={active}
              className={cn(
                "flex h-17 items-center justify-center rounded-control border px-3 text-center text-body text-primary transition-colors outline-none",
                "border-line bg-surface",
                active && "border-2 border-primary bg-highlight-soft",
                done && "border-2 border-highlight bg-highlight-soft text-muted",
                miss && "border-2 border-danger bg-danger-soft",
              )}
            >
              {item.text}
            </button>
          );
        })}
      </div>
    </div>
  );
}
