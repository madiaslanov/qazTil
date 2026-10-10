"use client";

import { useCallback, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { useLearnerStore } from "@/entities/learner";
import { progressQueries } from "@/entities/progress";
import { quizApi, quizQueries } from "@/entities/quiz";
import { ApiError } from "@/shared/api";

function messageOf(cause: unknown, fallback: string) {
  return cause instanceof ApiError ? cause.message : fallback;
}

/**
 * answering — выбираем вариант;
 * feedback — показан разбор ответа;
 * retrying — повтор после ошибки: сервер ответ уже засчитал,
 * поэтому вторую попытку проверяем локально по известному ключу.
 */
type Phase = "answering" | "feedback" | "retrying";

/**
 * Один проход урока: вопросы и счёт держит сервер,
 * хук отвечает за текущий вопрос и состояние проверки.
 */
export function useQuizSession(categoryId: number, size: number) {
  const queryClient = useQueryClient();
  const loseLife = useLearnerStore((state) => state.loseLife);
  const completeLesson = useLearnerStore((state) => state.completeLesson);

  const options = quizQueries.session(categoryId, size);
  const session = useQuery(options);

  const [index, setIndex] = useState(0);
  const [selected, setSelected] = useState<number | null>(null);
  const [correctIndex, setCorrectIndex] = useState<number | null>(null);
  const [phase, setPhase] = useState<Phase>("answering");
  const [finished, setFinished] = useState(false);

  const quiz = session.data ?? null;

  const answer = useMutation({
    mutationFn: (variables: { questionId: number; selectedIndex: number }) => {
      if (!quiz) throw new ApiError("урок ещё не готов", 0);
      return quizApi.answer(
        quiz.id,
        variables.questionId,
        variables.selectedIndex,
      );
    },
    onSuccess: (result) => {
      queryClient.setQueryData(options.queryKey, result.quiz);
      setCorrectIndex(result.correct_index);
      setPhase("feedback");
      if (!result.correct) loseLife();
    },
  });

  const question = quiz?.questions[index] ?? null;
  const total = quiz?.questions.length ?? 0;

  const check = useCallback(() => {
    if (!question || selected === null || phase === "feedback") return;
    if (phase === "retrying") {
      setPhase("feedback");
      return;
    }
    answer.mutate({ questionId: question.id, selectedIndex: selected });
  }, [answer, phase, question, selected]);

  const retry = useCallback(() => {
    setSelected(null);
    setPhase("retrying");
  }, []);

  const next = useCallback(() => {
    if (!quiz) return;
    setSelected(null);
    setCorrectIndex(null);
    setPhase("answering");

    if (index + 1 >= quiz.questions.length) {
      completeLesson(quiz.category_id, quiz.score.correct);
      void queryClient.invalidateQueries({
        queryKey: progressQueries.all().queryKey,
      });
      setFinished(true);
      return;
    }
    setIndex(index + 1);
  }, [completeLesson, index, queryClient, quiz]);

  const status = finished
    ? "finished"
    : session.isPending
      ? "loading"
      : session.isError || !question
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
    question,
    index,
    total,
    selected,
    /** Правильный вариант известен только после проверки. */
    correctIndex: checked ? correctIndex : null,
    checked,
    /** Верен ли последний проверенный ответ. */
    correct: checked ? selected === correctIndex : null,
    checking: answer.isPending,
    isLast: index + 1 >= total,
    select: setSelected,
    check,
    retry,
    next,
    /** Доля пройденных вопросов для полосы прогресса. */
    percent: total === 0 ? 0 : Math.round((index / total) * 100),
  };
}
