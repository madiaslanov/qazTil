"use client";

import Image from "next/image";
import Link from "next/link";
import { Check, Flame, Target } from "lucide-react";

import { useLearnerStore, XP_PER_CORRECT_ANSWER } from "@/entities/learner";
import type { Quiz } from "@/entities/quiz";
import { Button, Card, IconBadge, Screen } from "@/shared/ui";

function plural(days: number) {
  const tail = days % 10;
  if (days > 10 && days < 20) return "дней";
  if (tail === 1) return "день";
  if (tail >= 2 && tail <= 4) return "дня";
  return "дней";
}

/** Экран после урока: XP, страйк и точность прохождения. */
export function LessonResult({ quiz }: { quiz: Quiz }) {
  const streak = useLearnerStore((state) => state.learner?.streak ?? 0);
  const accuracy =
    quiz.score.total === 0
      ? 0
      : Math.round((quiz.score.correct / quiz.score.total) * 100);

  return (
    <Screen>
      <Image
        src="/bg/complete.png"
        alt=""
        fill
        sizes="430px"
        className="pointer-events-none object-cover mix-blend-color-burn"
      />

      <div className="relative flex flex-1 flex-col px-6 pt-[58px] pb-7">
        <div className="flex flex-col items-center gap-6">
          <span className="flex size-[104px] items-center justify-center rounded-full border-4 border-ink bg-green shadow-hard">
            <Check className="size-12" strokeWidth={3} />
          </span>

          <div className="flex flex-col items-center gap-2 text-center">
            <p className="text-[22px] text-muted">
              +{quiz.score.correct * XP_PER_CORRECT_ANSWER} XP
            </p>
            <h1 className="text-[38px] font-black leading-[1.02]">
              Урок завершен
            </h1>
            <p className="max-w-[292px] text-[15px] leading-[1.45] text-muted">
              Правильных ответов {quiz.score.correct} из {quiz.score.total}.
              Новые слова уже в твоём словаре.
            </p>
          </div>

          <Card className="flex w-full flex-col gap-4.5 p-5">
            <div className="flex items-center gap-3">
              <IconBadge className="bg-orange">
                <Flame strokeWidth={2.2} />
              </IconBadge>
              <div className="flex flex-col gap-0.5">
                <p className="text-[12px] font-bold uppercase text-subtle">
                  Страйк
                </p>
                <p className="text-[19px]">
                  {streak} {plural(streak)}
                </p>
              </div>
            </div>

            <div className="h-0.5 w-full bg-ink" />

            <div className="flex items-center gap-3">
              <IconBadge className="bg-green">
                <Target strokeWidth={2.2} />
              </IconBadge>
              <div className="flex flex-col gap-0.5">
                <p className="text-[12px] font-bold uppercase text-subtle">
                  Точность
                </p>
                <p className="text-[19px]">{accuracy}%</p>
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
