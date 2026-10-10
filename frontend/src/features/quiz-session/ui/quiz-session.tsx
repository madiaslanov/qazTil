"use client";

import { useEffect } from "react";
import Link from "next/link";

import type { Quiz } from "@/entities/quiz";
import { cn } from "@/shared/lib/cn";
import { Button, Eyebrow, FieldHint, StateNote } from "@/shared/ui";

import { useQuizSession } from "../model/use-quiz-session";
import { AnswerFeedback } from "./answer-feedback";
import { AnswerTile } from "./answer-tile";
import { LessonHeader } from "./lesson-header";

/** Экран прохождения урока: прогресс, вопрос, варианты и разбор ответа. */
export function QuizSession({
  categoryId,
  size,
  onFinished,
}: {
  categoryId: number;
  size: number;
  onFinished: (quiz: Quiz) => void;
}) {
  const session = useQuizSession(categoryId, size);
  const { status, quiz } = session;

  useEffect(() => {
    if (status === "finished" && quiz) {
      onFinished(quiz);
    }
  }, [onFinished, quiz, status]);

  if (session.status === "loading") {
    return <StateNote text="Собираем урок…" />;
  }

  if (session.status === "failed" || !session.question) {
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

  const { question } = session;

  return (
    <div className="relative flex min-h-0 flex-1 flex-col">
      <LessonHeader percent={session.percent} />

      <div className="relative flex min-h-0 flex-1 flex-col">
        <div
          className={cn(
            "flex flex-1 flex-col gap-6 overflow-y-auto px-gutter pt-8.5 pb-6 transition-opacity",
            session.checked && "opacity-68",
          )}
        >
          <div className="flex flex-col gap-2.25">
            <Eyebrow>Выбери правильный перевод</Eyebrow>
            <h1 className="text-h2 text-foreground">
              Как сказать «{question.prompt}»?
            </h1>
          </div>

          <div className="flex flex-col gap-3">
            {question.options.map((option, index) => (
              <AnswerTile
                key={`${question.id}-${index}`}
                index={index}
                option={option}
                selected={session.selected}
                correctIndex={session.correctIndex}
                disabled={session.checked || session.checking}
                onSelect={session.select}
              />
            ))}
          </div>

          {session.error && <FieldHint tone="error">{session.error}</FieldHint>}
        </div>

        {session.checked && (
          <div aria-hidden className="absolute inset-0 bg-scrim animate-fade-in" />
        )}
      </div>

      <footer className="shrink-0 bg-surface px-gutter pt-4.5 pb-7">
        <Button
          onClick={session.check}
          disabled={session.selected === null || session.checking || session.checked}
        >
          {session.checking ? "Проверяем…" : "Проверить ответ"}
        </Button>
      </footer>

      {session.checked && session.correctIndex !== null && (
        <AnswerFeedback
          correct={session.correct === true}
          answer={question.options[session.correctIndex]}
          isLast={session.isLast}
          onNext={session.next}
          onRetry={session.retry}
        />
      )}
    </div>
  );
}
