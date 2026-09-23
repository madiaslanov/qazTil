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
    "flex items-center justify-center rounded-full border-[3px] border-ink transition-transform",
    state === "current"
      ? "size-[76px] bg-orange shadow-hard active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
      : "size-[62px]",
    state === "completed" && "bg-ink text-white",
    state === "locked" && "bg-paper text-ink",
  );
  const content = (
    <>
      <span className={circle}>
        <Icon className="size-[26px]" strokeWidth={2.2} />
      </span>
      <span className="text-[13px] font-bold text-ink">{label}</span>
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
