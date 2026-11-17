"use client";

import { useCallback, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import type { Exercise, ExerciseAnswer, ExerciseResult } from "@/entities/exercise";
import { lessonApi, type Lesson } from "@/entities/lesson";
import { profileQueries } from "@/entities/learner";
import { pathQueries } from "@/entities/path";
import { progressQueries } from "@/entities/progress";
import { ApiError } from "@/shared/api";

/** answering — отвечаем; feedback — показан разбор ответа. */
type Phase = "answering" | "feedback";

/** Итог урока: то, что накопили ответы сервера за этот проход. */
export type LessonSummary = {
  xpEarned: number;
  streak: number;
  streakIncreased: boolean;
  nextLessonId: number | null;
  correct: number;
  total: number;
};

/** Уже сыгранная реплика сцены: что сказал собеседник и что ответил ученик. */
export type DialogueTurn = {
  speaker: string;
  line: string;
  hint: string;
  reply: string;
};

function messageOf(cause: unknown, fallback: string) {
  return cause instanceof ApiError ? cause.message : fallback;
}

/**
 * Один проход урока. Порядок, проверка, жизни, XP и стрик — всё на сервере:
 * хук только показывает ответы и ведёт по упражнениям с точки продолжения.
 * Урок открывается один раз на вход, поэтому состояние стартует от
 * resume_position и дальше живёт локально.
 */
export function useLessonSession(lesson: Lesson) {
  const queryClient = useQueryClient();
  const exercises = lesson.exercises;
  const total = exercises.length;

  const [index, setIndex] = useState(
    Math.min(lesson.resume_position, Math.max(0, total - 1)),
  );
  const [phase, setPhase] = useState<Phase>("answering");
  const [result, setResult] = useState<ExerciseResult | null>(null);
  const [selected, setSelected] = useState<number | null>(null);
  /** Растёт при повторе: сбрасывает состояние упражнения. */
  const [attempt, setAttempt] = useState(0);
  const [lives, setLives] = useState(lesson.lives);
  /** Реплики сцены, которые ученик уже прошёл в этом заходе. */
  const [turns, setTurns] = useState<DialogueTurn[]>([]);
  const [failed, setFailed] = useState(false);
  const [finished, setFinished] = useState(false);
  const [summary, setSummary] = useState<LessonSummary>({
    xpEarned: 0,
    streak: 0,
    streakIncreased: false,
    nextLessonId: null,
    // У пройденного урока счёт начинается заново.
    correct: lesson.state === "completed" ? 0 : lesson.correct,
    total,
  });

  const step: Exercise | null = exercises[index] ?? null;

  const answer = useMutation({
    mutationFn: (body: ExerciseAnswer) => {
      if (!step) throw new ApiError("упражнение не найдено", 0);
      return lessonApi.answer(step.id, body);
    },
    onSuccess: (response) => {
      setResult(response);
      setLives(response.lives);
      setPhase("feedback");
      setSummary((current) => ({
        ...current,
        xpEarned: current.xpEarned + response.xp_earned,
        streak: response.streak,
        streakIncreased: current.streakIncreased || response.streak_increased,
        nextLessonId: response.next_lesson_id ?? current.nextLessonId,
        // Первая попытка — attempt 0; повтор очков не добавляет.
        correct: current.correct + (response.correct && attempt === 0 ? 1 : 0),
      }));
      // Жизни и XP в шапке берутся из профиля.
      void queryClient.invalidateQueries({
        queryKey: profileQueries.current().queryKey,
      });
    },
    onError: (cause) => {
      // 409 — жизни кончились (в том числе в другой вкладке).
      if (cause instanceof ApiError && cause.status === 409) {
        setLives(0);
        setFailed(true);
        void queryClient.invalidateQueries({
          queryKey: profileQueries.current().queryKey,
        });
      }
    },
  });

  const submit = useCallback(
    (body: ExerciseAnswer) => {
      if (phase !== "answering" || answer.isPending) return;
      answer.mutate(body);
    },
    [answer, phase],
  );

  const retry = useCallback(() => {
    setSelected(null);
    setResult(null);
    setAttempt((value) => value + 1);
    setPhase("answering");
  }, []);

  const next = useCallback(() => {
    // Последняя жизнь ушла на этом ответе: урок остановлен.
    if (result?.lesson_failed) {
      setFailed(true);
      return;
    }
    // Реплика уходит в историю с верным ответом, который отдал сервер.
    if (step?.kind === "dialogue" && result) {
      setTurns((previous) => [
        ...previous,
        {
          speaker: step.speaker,
          line: step.prompt,
          hint: step.hint,
          // Ошибка присылает все подходящие ответы через « / »: в историю
          // сцены берём первый, это самый простой.
          reply: result.expected.split(" / ")[0],
        },
      ]);
    }
    setSelected(null);
    setResult(null);
    setAttempt(0);
    setPhase("answering");

    if (index + 1 >= total) {
      for (const key of [
        pathQueries.all().queryKey,
        progressQueries.all().queryKey,
        profileQueries.current().queryKey,
      ]) {
        void queryClient.invalidateQueries({ queryKey: key });
      }
      setFinished(true);
      return;
    }
    setIndex(index + 1);
  }, [index, queryClient, result, step, total]);

  const status = failed ? "failed" : finished ? "finished" : "running";
  const checked = phase === "feedback";

  /** Для выбора подсвечиваем верный вариант по тексту из expected. */
  const correctIndex =
    checked && step?.kind === "choice"
      ? step.options.findIndex((option) => option.text === result?.expected)
      : -1;

  return {
    status,
    error: answer.isError && !failed ? messageOf(answer.error, "ответ не сохранился") : null,
    step,
    turns,
    index,
    total,
    attempt,
    lives,
    maxLives: lesson.max_lives,
    selected,
    correctIndex: correctIndex >= 0 ? correctIndex : null,
    checked,
    result: checked ? result : null,
    checking: answer.isPending,
    isLast: index + 1 >= total,
    summary,
    select: setSelected,
    submit,
    retry,
    next,
    percent: total === 0 ? 0 : Math.round((index / total) * 100),
  };
}
