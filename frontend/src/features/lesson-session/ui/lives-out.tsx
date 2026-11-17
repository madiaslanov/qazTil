"use client";

import Link from "next/link";
import { Clock3, Heart, HeartCrack } from "lucide-react";

import {
  formatCountdown,
  LIFE_REGEN_MS,
  MAX_LIVES,
  useLives,
} from "@/entities/learner";
import { cn } from "@/shared/lib/cn";
import { Button } from "@/shared/ui";

/**
 * Урок остановлен: жизней нет. Показывает, когда вернётся следующая, и сам
 * оживает, как только профиль отдаёт хотя бы одну: тогда можно продолжить
 * с того места, где остановились.
 */
export function LivesOut({ onResume }: { onResume: () => void }) {
  const { lives = 0, msToNextLife } = useLives();
  const canResume = lives > 0;
  const minutes = LIFE_REGEN_MS / 60_000;

  return (
    <div className="flex flex-1 flex-col items-center gap-7 px-gutter pt-20 pb-7.5">
      <span className="flex size-29 items-center justify-center rounded-full bg-danger-soft">
        <HeartCrack className="size-15.5 text-danger" strokeWidth={1.5} />
      </span>

      <div className="flex flex-col gap-2 text-center">
        <h1 className="text-h1 font-bold text-heading">
          {canResume ? "Жизнь вернулась" : "Жизни закончились"}
        </h1>
        <p className="max-w-73 text-body-sm text-muted">
          {canResume
            ? "Можно продолжать: урок откроется с того упражнения, где ты остановился."
            : "Урок на паузе. Ответы сохранены — вернёшься с того же места."}
        </p>
      </div>

      <div
        className="flex gap-2"
        role="img"
        aria-label={`Жизни: ${lives} из ${MAX_LIVES}`}
      >
        {Array.from({ length: MAX_LIVES }, (_, index) => (
          <Heart
            key={index}
            className={cn(
              "size-8",
              index < lives ? "fill-danger text-danger" : "text-line",
            )}
          />
        ))}
      </div>

      {!canResume && msToNextLife !== null && (
        <div className="flex w-full items-center justify-between rounded-control bg-surface px-4 py-3.5">
          <span className="flex items-center gap-2 text-caption font-bold text-foreground">
            <Clock3 className="size-5 text-primary" />
            Новая жизнь через
          </span>
          <span className="text-stat font-extrabold text-danger tabular-nums">
            {formatCountdown(msToNextLife)}
          </span>
        </div>
      )}

      {!canResume && (
        <p className="text-center text-caption text-muted">
          Одна жизнь возвращается за {minutes} минут. Жизнь тратится только на
          первую ошибку в упражнении.
        </p>
      )}

      <div className="mt-auto flex w-full flex-col gap-2">
        <Button onClick={onResume} disabled={!canResume}>
          {canResume ? "Продолжить урок" : "Ждём новую жизнь"}
        </Button>
        <Button asChild variant="ghost" size="sm">
          <Link href="/learn">Вернуться на путь</Link>
        </Button>
      </div>
    </div>
  );
}
