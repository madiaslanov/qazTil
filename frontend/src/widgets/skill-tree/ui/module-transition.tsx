import { ArrowRight, Lock } from "lucide-react";

import { cn } from "@/shared/lib/cn";

/** Карточка между модулями: «Новый модуль открыт · A1 → A2». */
export function ModuleTransition({
  from,
  to,
  unlocked,
  onStart,
}: {
  from: string;
  to: string;
  unlocked: boolean;
  onStart: () => void;
}) {
  return (
    <div
      className={cn(
        "flex w-full items-center justify-between gap-4 rounded-card border-2 border-surface px-5 py-5.5",
        unlocked
          ? "bg-linear-to-b from-highlight-soft to-sky shadow-module"
          : "bg-surface-muted",
      )}
    >
      <div className="flex flex-col gap-1.25">
        <p className="text-micro font-extrabold text-accent uppercase">
          {unlocked ? "Новый модуль открыт" : "Следующий модуль"}
        </p>
        <p className="text-title text-primary">
          {from} → {to}
        </p>
        <p className="max-w-52.5 text-eyebrow text-muted">
          Переходим от базовых фраз к свободным диалогам.
        </p>
      </div>
      {unlocked ? (
        <button
          type="button"
          onClick={onStart}
          aria-label={`Начать модуль ${to}`}
          className="flex size-11 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground transition-transform active:scale-95"
        >
          <ArrowRight className="size-5" />
        </button>
      ) : (
        <span
          title="Пройди все уроки модуля"
          className="flex size-11 shrink-0 items-center justify-center rounded-full bg-line text-muted"
        >
          <Lock className="size-5" />
        </span>
      )}
    </div>
  );
}
