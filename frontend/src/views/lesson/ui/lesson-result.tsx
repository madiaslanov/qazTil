"use client";

import Image from "next/image";
import Link from "next/link";
import { Check, Flame, Target } from "lucide-react";

import { useLearnerStore, XP_PER_CORRECT_ANSWER } from "@/entities/learner";
import type { LessonScore } from "@/features/quiz-session";
import { Button, Card, IconBadge, Screen } from "@/shared/ui";

function plural(days: number) {
  const tail = days % 10;
  if (days > 10 && days < 20) return "дней";
  if (tail === 1) return "день";
  if (tail >= 2 && tail <= 4) return "дня";
  return "дней";
}

/** Экран после урока: XP, страйк и точность прохождения. */
export function LessonResult({ score }: { score: LessonScore }) {
  const streak = useLearnerStore((state) => state.learner?.streak ?? 0);
  const accuracy =
    score.total === 0
      ? 0
      : Math.round((score.correct / score.total) * 100);

  return (
    <Screen>
      <Image
        src="/bg/complete.png"
        alt=""
        fill
        sizes="430px"
        className="pointer-events-none object-cover mix-blend-color-burn"
      />

      <div className="relative flex flex-1 flex-col px-6 pt-14.5 pb-7">
        <div className="flex flex-col items-center gap-6">
          <span className="flex size-26 items-center justify-center rounded-full border-4 border-primary bg-success shadow-button">
            <Check className="size-12" strokeWidth={3} />
          </span>

          <div className="flex flex-col items-center gap-2 text-center">
            <p className="text-title text-muted">
              +{score.correct * XP_PER_CORRECT_ANSWER} XP
            </p>
            <h1 className="text-hero font-black">
              Урок завершен
            </h1>
            <p className="max-w-73 text-body text-muted">
              Правильных ответов {score.correct} из {score.total}.
              Новые слова уже в твоём словаре.
            </p>
          </div>

          <Card className="flex w-full flex-col gap-4.5 p-5">
            <div className="flex items-center gap-3">
              <IconBadge className="bg-accent">
                <Flame strokeWidth={2.2} />
              </IconBadge>
              <div className="flex flex-col gap-0.5">
                <p className="text-eyebrow font-bold uppercase text-muted">
                  Страйк
                </p>
                <p className="text-stat">
                  {streak} {plural(streak)}
                </p>
              </div>
            </div>

            <div className="h-0.5 w-full bg-primary" />

            <div className="flex items-center gap-3">
              <IconBadge className="bg-success">
                <Target strokeWidth={2.2} />
              </IconBadge>
              <div className="flex flex-col gap-0.5">
                <p className="text-eyebrow font-bold uppercase text-muted">
                  Точность
                </p>
                <p className="text-stat">{accuracy}%</p>
              </div>
            </div>
          </Card>
        </div>

        <div className="mt-auto pt-8">
          <Button asChild variant="primary">
            <Link href="/learn">Продолжить</Link>
          </Button>
        </div>
      </div>
    </Screen>
  );
}
