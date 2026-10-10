"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

import type { AssembleExercise, MatchExercise } from "@/entities/exercise";
import { cn } from "@/shared/lib/cn";
import { Button, Eyebrow, FieldHint, StateNote } from "@/shared/ui";

import { useQuizSession, type Step } from "../model/use-quiz-session";
import { AnswerFeedback } from "./answer-feedback";
import { AnswerTile } from "./answer-tile";
import { LessonHeader } from "./lesson-header";
import { MatchBoard } from "./match-board";
import { PhraseBuilder } from "./phrase-builder";

export type LessonScore = { correct: number; total: number };

function headingOf(step: Step) {
  switch (step.kind) {
    case "choice":
      return {
        eyebrow: "Выбери правильный перевод",
        title: `Как сказать «${step.question.prompt}»?`,
      };
    case "assemble":
      return {
        eyebrow: "Собери фразу",
        title: `«${(step.exercise as AssembleExercise).prompt}»`,
      };
    case "match":
      return { eyebrow: "Найди пару", title: "Соедини слово и перевод" };
  }
}

/** Экран прохождения урока: прогресс, задание, ответ и его разбор. */
export function QuizSession({
  categoryId,
  size,
  onFinished,
}: {
  categoryId: number;
  size: number;
  onFinished: (score: LessonScore) => void;
}) {
  const session = useQuizSession(categoryId, size);
  const { status, score } = session;
  // Собранная фраза живёт, пока не сменился шаг или попытка.
  const stepKey = `${session.index}-${session.attempt}`;
  const [assembled, setAssembled] = useState({ key: "", ids: [] as string[] });
  const tileIds = assembled.key === stepKey ? assembled.ids : [];

  useEffect(() => {
    if (status === "finished") onFinished(score);
  }, [onFinished, score, status]);

  if (session.status === "loading") {
    return <StateNote text="Собираем урок…" />;
  }

  if (session.status === "failed" || !session.step) {
    return (
      <StateNote
        text={session.error ?? "урок не собрался"}
        action={
          <Button asChild size="md" variant="secondary" className="w-auto">
            <Link href="/learn">Вернуться к пути</Link>
          </Button>
        }
      />
    );
  }

  const { step } = session;
  const heading = headingOf(step);
  const locked = session.checked || session.checking;

  return (
    <div className="relative flex min-h-0 flex-1 flex-col">
      <LessonHeader percent={session.percent} />

      <div className="relative flex min-h-0 flex-1 flex-col">
        <div
          className={cn(
            "flex flex-1 flex-col gap-7 overflow-y-auto px-gutter pt-8.5 pb-6 transition-opacity",
            session.checked && "opacity-68",
          )}
        >
          <div className="flex flex-col gap-2.25">
            <Eyebrow>{heading.eyebrow}</Eyebrow>
            <h1 className="text-h2 text-foreground">{heading.title}</h1>
          </div>

          {step.kind === "choice" && (
            <div className="flex flex-col gap-3">
              {step.question.options.map((option, index) => (
                <AnswerTile
                  key={`${step.question.id}-${index}`}
                  index={index}
                  option={option}
                  selected={session.selected}
                  correctIndex={session.correctIndex}
                  disabled={locked}
                  onSelect={session.select}
                />
              ))}
            </div>
          )}

          {step.kind === "assemble" && (
            <PhraseBuilder
              exercise={step.exercise as AssembleExercise}
              value={tileIds}
              disabled={locked}
              onChange={(update) =>
                setAssembled((current) => ({
                  key: stepKey,
                  ids: update(current.key === stepKey ? current.ids : []),
                }))
              }
            />
          )}

          {step.kind === "match" && (
            <MatchBoard
              key={stepKey}
              exercise={step.exercise as MatchExercise}
              disabled={locked}
              onComplete={(pairs) =>
                session.submitExercise({ kind: "match", pairs })
              }
            />
          )}

          {session.error && <FieldHint tone="error">{session.error}</FieldHint>}
        </div>

        {session.checked && (
          <div aria-hidden className="absolute inset-0 bg-scrim animate-fade-in" />
        )}
      </div>

      <footer className="shrink-0 bg-surface px-gutter pt-4.5 pb-7">
        {step.kind === "choice" && (
          <Button
            onClick={session.checkChoice}
            disabled={session.selected === null || locked}
          >
            {session.checking ? "Проверяем…" : "Проверить ответ"}
          </Button>
        )}
        {step.kind === "assemble" && (
          <Button
            onClick={() =>
              session.submitExercise({ kind: "assemble", tile_ids: tileIds })
            }
            disabled={tileIds.length === 0 || locked}
          >
            Проверить ответ
          </Button>
        )}
        {step.kind === "match" && (
          <Button disabled>Найди все пары</Button>
        )}
      </footer>

      {session.result && (
        <AnswerFeedback
          correct={session.result.correct}
          answer={session.result.answer}
          isLast={session.isLast}
          speakable={step.kind !== "match"}
          onNext={session.next}
          onRetry={session.retry}
        />
      )}
    </div>
  );
}
