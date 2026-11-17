"use client";

import Link from "next/link";
import { Check, Flame, Target } from "lucide-react";

import { useQuery } from "@tanstack/react-query";

import { profileQueries } from "@/entities/learner";
import type { LessonSummary } from "@/features/lesson-session";
import { Button, Card, IconBadge, Screen } from "@/shared/ui";

function plural(days: number) {
  const tail = days % 10;
  if (days % 100 > 10 && days % 100 < 20) return "дней";
  if (tail === 1) return "день";
  if (tail >= 2 && tail <= 4) return "дня";
  return "дней";
}

/** Экран после урока. XP, стрик и следующий урок — из ответов сервера. */
export function LessonResult({ summary }: { summary: LessonSummary }) {
  const profile = useQuery(profileQueries.current());
  const goalMet = profile.data?.today.goal_met ?? false;
  const { streak } = summary;
  const accuracy =
    summary.total === 0
      ? 0
      : Math.round((summary.correct / summary.total) * 100);

  return (
    <Screen className="px-gutter pt-22 pb-7.5">
      <div className="flex flex-col items-center gap-6">
        <span className="flex size-26 items-center justify-center rounded-full bg-highlight-strong shadow-glow">
          <Check className="size-12 text-black" />
        </span>

        <div className="flex flex-col items-center gap-2 text-center">
          <p className="text-title text-muted">
            +{summary.xpEarned} XP
          </p>
          <h1 className="text-hero font-black text-foreground">Урок завершен</h1>
          <p className="max-w-73 text-body text-muted">
            Верно {summary.correct} из {summary.total}.
            {goalMet && " Дневная цель выполнена."}
          </p>
        </div>

        <Card className="flex w-full flex-col gap-4.5 p-5">
          <Stat
            icon={<Flame />}
            tone="bg-lives"
            label={summary.streakIncreased ? "Страйк вырос" : "Страйк"}
            value={`${streak} ${plural(streak)}`}
          />
          <Stat
            icon={<Target />}
            tone="bg-highlight"
            label="Точность"
            value={`${accuracy}%`}
          />
        </Card>
      </div>

      <Button asChild className="mt-auto">
        <Link href="/learn">
          {summary.nextLessonId === null ? "К пути" : "Продолжить"}
        </Link>
      </Button>
    </Screen>
  );
}

function Stat({
  icon,
  tone,
  label,
  value,
}: {
  icon: React.ReactNode;
  tone: string;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-center gap-3">
      <IconBadge className={`${tone} text-foreground`}>{icon}</IconBadge>
      <div className="flex flex-col gap-0.5">
        <p className="text-eyebrow font-bold text-muted uppercase">{label}</p>
        <p className="text-stat text-foreground">{value}</p>
      </div>
    </div>
  );
}
