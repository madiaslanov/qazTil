import { ArrowRight, Lock } from "lucide-react";

import { cn } from "@/shared/lib/cn";

export type UnitBannerState = "open" | "completed" | "locked";

const eyebrows: Record<UnitBannerState, string> = {
  open: "открыт",
  completed: "пройден",
  locked: "закрыт",
};

/** Карточка в начале юнита: «Юнит 2 открыт · В магазине» со стрелкой на старт. */
export function UnitBanner({
  number,
  title,
  description,
  state,
  onStart,
}: {
  number: number;
  title: string;
  description: string;
  state: UnitBannerState;
  onStart: () => void;
}) {
  return (
    <div
      className={cn(
        "flex w-full items-center justify-between gap-4 rounded-card border-2 border-surface px-5 py-5.5",
        state === "open"
          ? "bg-linear-to-b from-highlight-soft to-sky shadow-module"
          : "bg-surface-muted",
      )}
    >
      <div className="flex flex-col gap-1.25">
        <p
          className={cn(
            "text-micro font-extrabold uppercase",
            state === "open" ? "text-accent" : "text-muted",
          )}
        >
          Юнит {number} {eyebrows[state]}
        </p>
        <p className="text-title text-primary">{title}</p>
        <p className="max-w-52.5 text-eyebrow text-muted">{description}</p>
      </div>
      {state === "locked" ? (
        <span
          title="Пройди предыдущий юнит"
          className="flex size-11 shrink-0 items-center justify-center rounded-full bg-line text-muted"
        >
          <Lock className="size-5" />
        </span>
      ) : (
        <button
          type="button"
          onClick={onStart}
          aria-label={`Открыть юнит «${title}»`}
          className="flex size-11 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground transition-transform active:scale-95"
        >
          <ArrowRight className="size-5" />
        </button>
      )}
    </div>
  );
}
