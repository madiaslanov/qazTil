import { CircleCheck, Lock, Play } from "lucide-react";

import { cn } from "@/shared/lib/cn";

export type SkillState = "completed" | "current" | "locked";

const icons = {
  completed: CircleCheck,
  current: Play,
  locked: Lock,
} as const;

const circles: Record<SkillState, string> = {
  completed: "size-15.5 bg-primary text-primary-foreground",
  current: "size-19 bg-highlight-strong text-foreground shadow-glow",
  locked: "size-15.5 bg-glow text-black",
};

/** Узел пути: пройдено, текущий урок или закрыто. Закрытый не нажимается. */
export function SkillNode({
  label,
  state,
  onSelect,
  className,
}: {
  label: string;
  state: SkillState;
  onSelect: () => void;
  className?: string;
}) {
  const Icon = icons[state];
  const locked = state === "locked";

  return (
    <button
      type="button"
      onClick={onSelect}
      disabled={locked}
      aria-current={state === "current" ? "step" : undefined}
      title={locked ? "Пройди предыдущий урок" : undefined}
      className={cn(
        "group flex flex-col items-center gap-2 outline-none disabled:cursor-not-allowed",
        state === "current" && "py-2.25",
        className,
      )}
    >
      <span
        className={cn(
          "flex items-center justify-center rounded-full transition-transform group-active:scale-95 group-disabled:group-active:scale-100",
          circles[state],
        )}
      >
        <Icon className="size-6.5" />
      </span>
      <span className="text-caption font-bold whitespace-nowrap text-foreground">
        {label}
      </span>
    </button>
  );
}
