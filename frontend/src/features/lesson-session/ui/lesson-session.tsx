"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

import type { Exercise, MatchExercise } from "@/entities/exercise";
import type { Lesson } from "@/entities/lesson";
import { cn } from "@/shared/lib/cn";
import { Button, Eyebrow, FieldHint, Input, StateNote } from "@/shared/ui";

import {
  useLessonSession,
  type LessonSummary,
} from "../model/use-lesson-session";
import { AnswerFeedback } from "./answer-feedback";
import { AnswerTile } from "./answer-tile";
import { DialogueChat, DialogueHelp } from "./dialogue-chat";
import { LessonHeader } from "./lesson-header";
import { LivesOut } from "./lives-out";
import { MatchBoard } from "./match-board";
import { PhraseBuilder } from "./phrase-builder";

function headingOf(step: Exercise) {
  switch (step.kind) {
    case "choice":
      return {
        eyebrow: "Выбери правильный перевод",
        title: `Как сказать «${step.prompt}»?`,
      };
    case "translate":
      return {
        eyebrow: "Переведи на казахский",
        title: `«${step.prompt}»`,
      };
    case "dialogue":
      return { eyebrow: "Ваша цель", title: step.goal };
    // У этих двух типов задание целиком в prompt: «Соберите: «…»».
    case "assemble":
      return { eyebrow: "Собери фразу", title: step.prompt };
    case "match":
      return {
        eyebrow: "Найди пару",
        title: step.prompt ?? "Соедини фразу с переводом",
      };
  }
}

type Draft = {
  key: string;
  tiles: string[];
  pairs: [string, string][];
  text: string;
};

/** Экран прохождения урока: прогресс, задание, ответ и его разбор. */
export function LessonSession({
  lesson,
  onFinished,
  onResume,
}: {
  lesson: Lesson;
  onFinished: (summary: LessonSummary) => void;
  /** Жизнь вернулась, ученик жмёт «Продолжить»: страница откроет урок заново. */
  onResume: () => void;
}) {
  const session = useLessonSession(lesson);
  const { status, summary } = session;

  // Черновик ответа живёт, пока не сменился шаг или попытка.
  const stepKey = `${session.index}-${session.attempt}`;
  const [stored, setStored] = useState<Draft>({
    key: "",
    tiles: [],
    pairs: [],
    text: "",
  });
  const draft: Draft =
    stored.key === stepKey
      ? stored
      : { key: stepKey, tiles: [], pairs: [], text: "" };
  const patch = (update: (current: Draft) => Partial<Draft>) =>
    setStored((previous) => {
      const base =
        previous.key === stepKey
          ? previous
          : { key: stepKey, tiles: [], pairs: [], text: "" };
      return { ...base, ...update(base) };
    });

  useEffect(() => {
    if (status === "finished") onFinished(summary);
  }, [onFinished, summary, status]);

  if (status === "failed") {
    return <LivesOut onResume={onResume} />;
  }

  const { step } = session;
  if (!step) {
    return (
      <StateNote
        text="В этом уроке пока нет упражнений"
        action={
          <Button asChild size="md" variant="secondary" className="w-auto">
            <Link href="/learn">Вернуться к пути</Link>
          </Button>
        }
      />
    );
  }

  const heading = headingOf(step);
  const locked = session.checked || session.checking;

  function check() {
    if (!step) return;
    switch (step.kind) {
      case "choice": {
        if (session.selected === null) return;
        const option = step.options[session.selected];
        session.submit({
          kind: "choice",
          option_id: option.id,
          answer: option.text,
        });
        return;
      }
      case "dialogue":
        session.submit({ kind: "dialogue", tile_ids: draft.tiles });
        return;
      case "translate":
        session.submit({ kind: "translate", answer: draft.text.trim() });
        return;
      case "assemble":
        session.submit({ kind: "assemble", tile_ids: draft.tiles });
        return;
      case "match":
        session.submit({ kind: "match", pairs: draft.pairs });
        return;
    }
  }

  const ready =
    step.kind === "choice"
      ? session.selected !== null
      : step.kind === "translate"
        ? draft.text.trim() !== ""
        : step.kind === "assemble" || step.kind === "dialogue"
          ? draft.tiles.length > 0
          : draft.pairs.length === step.pairs_left.length;

  const failedNow = session.result?.lesson_failed ?? false;

  return (
    <div className="relative flex min-h-0 flex-1 flex-col">
      <LessonHeader
        percent={session.percent}
        lives={session.lives}
        maxLives={session.maxLives}
      />

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

          {step.kind === "dialogue" && (
            <>
              <DialogueChat turns={session.turns} step={step} />
              {step.help && <DialogueHelp key={step.id} text={step.help} />}
            </>
          )}

          {step.kind === "choice" && (
            <div className="flex flex-col gap-3">
              {step.options.map((option, index) => (
                <AnswerTile
                  key={`${step.id}-${option.id}`}
                  index={index}
                  option={option.text}
                  selected={session.selected}
                  correctIndex={session.correctIndex}
                  disabled={locked}
                  onSelect={session.select}
                />
              ))}
            </div>
          )}

          {step.kind === "translate" && (
            <form
              onSubmit={(event) => {
                event.preventDefault();
                if (ready && !locked) check();
              }}
            >
              <Input
                key={stepKey}
                lang="kk"
                value={draft.text}
                onChange={(event) => patch(() => ({ text: event.target.value }))}
                placeholder="Напиши по-казахски"
                autoComplete="off"
                autoCapitalize="off"
                spellCheck={false}
                disabled={locked}
                aria-label="Перевод"
              />
            </form>
          )}

          {(step.kind === "assemble" || step.kind === "dialogue") && (
            <PhraseBuilder
              exercise={step}
              value={draft.tiles}
              disabled={locked}
              onChange={(update) =>
                patch((current) => ({ tiles: update(current.tiles) }))
              }
            />
          )}

          {step.kind === "match" && (
            <MatchBoard
              key={stepKey}
              exercise={step as MatchExercise}
              pairs={draft.pairs}
              disabled={locked}
              onChange={(pairs) => patch(() => ({ pairs }))}
            />
          )}

          {session.error && <FieldHint tone="error">{session.error}</FieldHint>}
        </div>

        {session.checked && (
          <div aria-hidden className="absolute inset-0 bg-scrim animate-fade-in" />
        )}
      </div>

      <footer className="shrink-0 bg-surface px-gutter pt-4.5 pb-7">
        <Button onClick={check} disabled={!ready || locked}>
          {session.checking
            ? "Проверяем…"
            : step.kind === "match" && !ready
              ? "Соедини все пары"
              : "Проверить ответ"}
        </Button>
      </footer>

      {session.result && (
        <AnswerFeedback
          correct={session.result.correct}
          answer={session.result.expected}
          isLast={session.isLast && !failedNow}
          speakable={step.kind !== "match"}
          canRetry={!failedNow}
          onNext={session.next}
          onRetry={session.retry}
        />
      )}
    </div>
  );
}
