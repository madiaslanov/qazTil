"use client";

import { Check, Volume2, X } from "lucide-react";

import { cn } from "@/shared/lib/cn";
import { canSpeak, speak } from "@/shared/lib/speech";
import { Button } from "@/shared/ui";

const copy = {
  correct: {
    title: "Дұрыс! Верно",
    text: "Фраза звучит естественно — так держать.",
    note: "Правильный ответ",
  },
  wrong: {
    title: "Почти получилось",
    text: "Сравни с правильным вариантом и попробуй ещё раз.",
    note: "Запомни",
  },
} as const;

/** Панель разбора снизу: «Верно» мятная, «Ошибка» коралловая. */
export function AnswerFeedback({
  correct,
  answer,
  isLast,
  onNext,
  onRetry,
}: {
  correct: boolean;
  /** Правильный вариант ответа. */
  answer: string;
  isLast: boolean;
  onNext: () => void;
  onRetry: () => void;
}) {
  const text = copy[correct ? "correct" : "wrong"];
  const Icon = correct ? Check : X;

  return (
    <section
      role="status"
      aria-live="polite"
      className={cn(
        "absolute inset-x-0 bottom-0 z-10 flex flex-col gap-4.5 rounded-t-sheet px-gutter pt-3 pb-7 animate-sheet-in",
        correct
          ? "bg-success-surface shadow-success"
          : "bg-danger-surface shadow-danger",
      )}
    >
      <div aria-hidden className="flex justify-center">
        <span
          className={cn(
            "h-1.25 w-11 rounded-full opacity-45",
            correct ? "bg-success" : "bg-danger",
          )}
        />
      </div>

      <div className="flex items-center gap-3.5">
        <span
          className={cn(
            "flex size-13.5 shrink-0 items-center justify-center rounded-full text-white",
            correct ? "bg-success" : "bg-danger",
          )}
        >
          <Icon className="size-6.25" />
        </span>
        <div className="flex flex-col gap-1">
          <p className="text-title text-heading">{text.title}</p>
          <p className="text-caption text-muted">{text.text}</p>
        </div>
      </div>

      <div className="flex items-center justify-between gap-3 rounded-control bg-surface/65 px-3.75 py-3.25">
        <div className="flex flex-col gap-0.75">
          <p
            className={cn(
              "text-nano font-extrabold uppercase",
              correct ? "text-success" : "text-danger",
            )}
          >
            {text.note}
          </p>
          <p lang="kk" className="text-lead text-primary">
            {answer}
          </p>
        </div>
        {canSpeak() && (
          <button
            type="button"
            onClick={() => speak(answer)}
            aria-label="Послушать"
            className="text-primary"
          >
            <Volume2 className="size-5.25" />
          </button>
        )}
      </div>

      {correct ? (
        <Button onClick={onNext}>
          {isLast ? "Завершить урок" : "Продолжить"}
        </Button>
      ) : (
        <>
          <Button onClick={onRetry}>Повторить</Button>
          <Button variant="ghost" size="sm" onClick={onNext} className="self-center text-primary">
            Продолжить без повтора
          </Button>
        </>
      )}
    </section>
  );
}
