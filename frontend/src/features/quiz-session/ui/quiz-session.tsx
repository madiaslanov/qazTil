"use client";

import { useEffect } from "react";
import Link from "next/link";
import { CircleX } from "lucide-react";

import type { Quiz } from "@/entities/quiz";
import { Button, Progress, StateNote } from "@/shared/ui";

import { useQuizSession } from "../model/use-quiz-session";
import { AnswerTile } from "./answer-tile";

/** Экран прохождения урока: полоса прогресса, вопрос, варианты, кнопка проверки. */
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
    <div className="flex min-h-0 flex-1 flex-col">
      <div className="flex h-bar shrink-0 items-center gap-4 border-b-3 border-primary px-5.5">
        <Link href="/learn" aria-label="Выйти из урока">
          <CircleX className="size-5.5" strokeWidth={2.2} />
        </Link>
        <Progress value={session.percent} className="flex-1" />
        <b className="text-caption font-extrabold">{session.percent}%</b>
      </div>

      <div className="flex flex-1 flex-col gap-7 overflow-y-auto px-6 pt-8.5 pb-6">
        <div className="flex flex-col gap-2.5">
          <p className="text-caption font-extrabold uppercase text-accent">
            Выбери правильный перевод
          </p>
          <h1 className="text-h2 text-foreground">
            Как сказать «{question.prompt}»?
          </h1>
        </div>

        <div className="flex flex-col gap-3.5">
          {question.options.map((option, index) => (
            <AnswerTile
              key={`${question.id}-${index}`}
              index={index}
              option={option}
              selected={session.selected}
              correctIndex={session.correctIndex}
              onSelect={session.select}
            />
          ))}
        </div>

        {session.error && (
          <p className="text-caption font-bold text-danger">{session.error}</p>
        )}
      </div>

      <div className="shrink-0 border-t-3 border-primary bg-surface px-6 pt-4.5 pb-7">
        {session.checked ? (
          <Button variant="primary" onClick={session.next}>
            {session.index + 1 >= session.total ? "Завершить урок" : "Продолжить"}
          </Button>
        ) : (
          <Button
            variant="primary"
            onClick={session.check}
            disabled={session.selected === null || session.checking}
          >
            Проверить ответ
          </Button>
        )}
      </div>
    </div>
  );
}
