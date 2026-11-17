import Link from "next/link";
import { Heart, X } from "lucide-react";

import { Progress } from "@/shared/ui";

/** Шапка урока: выход, полоса прогресса и остаток жизней. */
export function LessonHeader({
  percent,
  lives,
  maxLives,
}: {
  percent: number;
  lives: number;
  maxLives: number;
}) {
  return (
    <header className="flex h-bar shrink-0 items-center gap-4.5 bg-surface px-gutter">
      <Link href="/learn" aria-label="Выйти из урока" className="text-foreground">
        <X className="size-5" />
      </Link>
      <Progress value={percent} className="flex-1" aria-label="Пройдено упражнений" />
      <b
        className="flex items-center gap-1 text-caption font-extrabold text-foreground tabular-nums"
        aria-label={`Жизни: ${lives} из ${maxLives}`}
      >
        <Heart className="size-4.5 fill-current text-danger" />
        {lives}
      </b>
    </header>
  );
}
