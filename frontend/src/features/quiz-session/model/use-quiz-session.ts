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
      if (!result.correct) loseLife();
    },
  });

  const question = quiz?.questions[index] ?? null;
  const total = quiz?.questions.length ?? 0;

  const check = useCallback(() => {
    if (!question || selected === null || correctIndex !== null) return;
    answer.mutate({ questionId: question.id, selectedIndex: selected });
  }, [answer, correctIndex, question, selected]);

  const next = useCallback(() => {
    if (!quiz) return;
    setSelected(null);
    setCorrectIndex(null);

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
    correctIndex,
    checked: correctIndex !== null,
    checking: answer.isPending,
    select: setSelected,
    check,
    next,
    /** Доля пройденных вопросов для полосы прогресса. */
    percent: total === 0 ? 0 : Math.round((index / total) * 100),
  };
}
