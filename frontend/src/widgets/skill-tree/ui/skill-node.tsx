import Link from "next/link";
import { CircleCheck, Lock, Play } from "lucide-react";

import { cn } from "@/shared/lib/cn";

export type SkillState = "completed" | "current" | "locked";

const icons = {
  completed: CircleCheck,
  current: Play,
  locked: Lock,
} as const;

/** Узел пути обучения: пройдено, текущий урок или закрыто. */
export function SkillNode({
  href,
  label,
  state,
}: {
  href: string;
  label: string;
  state: SkillState;
}) {
  const Icon = icons[state];
  const circle = cn(
    "flex items-center justify-center rounded-full border-3 border-primary transition-transform",
    state === "current"
      ? "size-19 bg-accent shadow-button active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
      : "size-15.5",
    state === "completed" && "bg-primary text-white",
    state === "locked" && "bg-background text-foreground",
  );
  const content = (
    <>
      <span className={circle}>
        <Icon className="size-6.5" strokeWidth={2.2} />
      </span>
      <span className="text-caption font-bold text-foreground">{label}</span>
    </>
  );

  if (state === "locked") {
    return (
      <div
        className="flex flex-col items-center gap-2 opacity-70"
        aria-disabled
        title="Пройди предыдущий урок"
      >
        {content}
      </div>
    );
  }

  return (
    <Link href={href} className="flex flex-col items-center gap-2">
      {content}
    </Link>
  );
}
