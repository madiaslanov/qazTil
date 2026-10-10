import Link from "next/link";
import { X } from "lucide-react";

import { Progress } from "@/shared/ui";

/** Шапка урока: выход, полоса прогресса и процент. */
export function LessonHeader({ percent }: { percent: number }) {
  return (
    <header className="flex h-bar shrink-0 items-center gap-4.5 bg-surface px-gutter">
      <Link href="/learn" aria-label="Выйти из урока" className="text-foreground">
        <X className="size-5" />
      </Link>
      <Progress value={percent} className="flex-1" aria-label="Пройдено вопросов" />
      <b className="w-9 text-right text-caption font-extrabold text-foreground tabular-nums">
        {percent}%
      </b>
    </header>
  );
}
