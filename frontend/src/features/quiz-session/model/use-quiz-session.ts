"use client";

import { useCallback, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  buildLocalExercises,
  checkLocalExercise,
  type Exercise,
  type ExerciseAnswer,
} from "@/entities/exercise";
import { useLearnerStore } from "@/entities/learner";
import { progressQueries } from "@/entities/progress";
import { quizApi, quizQueries, type Question } from "@/entities/quiz";
import { wordQueries } from "@/entities/word";
import { ApiError } from "@/shared/api";

function messageOf(cause: unknown, fallback: string) {
  return cause instanceof ApiError ? cause.message : fallback;
}

export type Step =
  | { kind: "choice"; question: Question }
  | { kind: "match" | "assemble"; exercise: Exercise };

/** Куда вставлять упражнения из слов: после второго и четвёртого вопроса. */
const EXERCISE_SLOTS = [2, 4];

function buildSteps(questions: Question[], exercises: Exercise[]): Step[] {
  const steps: Step[] = questions.map((question) => ({ kind: "choice", question }));
  exercises.forEach((exercise, index) => {
    const at = Math.min(EXERCISE_SLOTS[index] + index, steps.length);
    steps.splice(at, 0, { kind: exercise.kind, exercise });
  });
  return steps;
}

/**
 * answering — отвечаем; feedback — показан разбор;
 * retrying — повтор после ошибки: первая попытка уже засчитана,
 * вторую проверяем без сервера и без потери жизни.
 */
type Phase = "answering" | "feedback" | "retrying";

type Result = { correct: boolean; answer: string };

/**
 * Один проход урока. Вопросы с выбором и счёт по ним держит сервер,
 * упражнения из слов проверяет локальный адаптер.
 */
export function useQuizSession(categoryId: number, size: number) {
  const queryClient = useQueryClient();
  const loseLife = useLearnerStore((state) => state.loseLife);
  const completeLesson = useLearnerStore((state) => state.completeLesson);

  const options = quizQueries.session(categoryId, size);
  const session = useQuery(options);
  const words = useQuery(wordQueries.list(categoryId));

  const [index, setIndex] = useState(0);
  const [phase, setPhase] = useState<Phase>("answering");
  const [result, setResult] = useState<Result | null>(null);
  const [selected, setSelected] = useState<number | null>(null);
  const [correctIndex, setCorrectIndex] = useState<number | null>(null);
  /** Меняется при повторе — сбрасывает состояние упражнения. */
  const [attempt, setAttempt] = useState(0);
  const [localScore, setLocalScore] = useState({ correct: 0, total: 0 });
  const [finished, setFinished] = useState(false);

  const quiz = session.data ?? null;
  // Состав урока фиксируем один раз, иначе перемешивание сбивало бы шаги.
  const steps = useMemo(
    () =>
      quiz && words.data
        ? buildSteps(quiz.questions, buildLocalExercises(words.data))
        : [],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [quiz?.id, words.data],
  );
  const step = steps[index] ?? null;
  const total = steps.length;

  const answer = useMutation({
    mutationFn: (variables: { questionId: number; selectedIndex: number }) => {
      if (!quiz) throw new ApiError("урок ещё не готов", 0);
      return quizApi.answer(quiz.id, variables.questionId, variables.selectedIndex);
    },
    onSuccess: (response, variables) => {
      queryClient.setQueryData(options.queryKey, response.quiz);
      setCorrectIndex(response.correct_index);
      const question = response.quiz.questions.find((q) => q.id === variables.questionId);
      setResult({
        correct: response.correct,
        answer: question?.options[response.correct_index] ?? "",
      });
      setPhase("feedback");
      if (!response.correct) loseLife();
    },
  });

  const checkChoice = useCallback(() => {
    if (step?.kind !== "choice" || selected === null) return;
    if (phase === "retrying") {
      setResult((previous) => previous && { ...previous, correct: selected === correctIndex });
      setPhase("feedback");
      return;
    }
    answer.mutate({ questionId: step.question.id, selectedIndex: selected });
  }, [answer, correctIndex, phase, selected, step]);

  /** Ответ на упражнение из слов: проверка локально, как потом на сервере. */
  const submitExercise = useCallback(
    (response: ExerciseAnswer) => {
      if (!step || step.kind === "choice" || phase === "feedback") return;
      const checked = checkLocalExercise(step.exercise, response);
      setResult({ correct: checked.correct, answer: checked.correct_answer });
      setPhase("feedback");
      if (phase === "answering") {
        setLocalScore((score) => ({
          correct: score.correct + (checked.correct ? 1 : 0),
          total: score.total + 1,
        }));
        if (!checked.correct) loseLife();
      }
    },
    [loseLife, phase, step],
  );

  const retry = useCallback(() => {
    setSelected(null);
    setAttempt((value) => value + 1);
    setPhase("retrying");
  }, []);

  const next = useCallback(() => {
    if (!quiz) return;
    setSelected(null);
    setCorrectIndex(null);
    setResult(null);
    setAttempt(0);
    setPhase("answering");

    if (index + 1 >= steps.length) {
      completeLesson(quiz.category_id, quiz.score.correct + localScore.correct);
      void queryClient.invalidateQueries({
        queryKey: progressQueries.all().queryKey,
      });
      setFinished(true);
      return;
    }
    setIndex(index + 1);
  }, [completeLesson, index, localScore.correct, queryClient, quiz, steps.length]);

  const status = finished
    ? "finished"
    : session.isPending || words.isPending
      ? "loading"
      : session.isError || !step
        ? "failed"
        : "running";

  const checked = phase === "feedback";

  return {
    status,
    error: session.isError
      ? messageOf(session.error, "не вышло начать урок")
      : answer.isError
        ? messageOf(answer.error, "ответ не сохранился")
        : null,
    quiz,
    step,
    index,
    total,
    attempt,
    selected,
    /** Правильный вариант известен только после проверки. */
    correctIndex: checked ? correctIndex : null,
    checked,
    result: checked ? result : null,
    checking: answer.isPending,
    isLast: index + 1 >= total,
    /** Итог урока: вопросы сервера плюс упражнения из слов. */
    score: {
      correct: (quiz?.score.correct ?? 0) + localScore.correct,
      total: (quiz?.questions.length ?? 0) + localScore.total,
    },
    select: setSelected,
    checkChoice,
    submitExercise,
    retry,
    next,
    /** Доля пройденных шагов для полосы прогресса. */
    percent: total === 0 ? 0 : Math.round((index / total) * 100),
  };
}
