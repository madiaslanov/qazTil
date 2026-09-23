"use client";

import { useCallback, useEffect, useState } from "react";

import { learnerStore } from "@/entities/learner";
import { quizApi, type Quiz } from "@/entities/quiz";
import { ApiError } from "@/shared/api";

type Status = "loading" | "running" | "finished" | "failed";

/**
 * Один проход квиза: сервер держит вопросы и счёт,
 * хук — только текущий вопрос и состояние проверки.
 */
export function useQuizSession(categoryId: number, size: number) {
  const [status, setStatus] = useState<Status>("loading");
  const [error, setError] = useState<string | null>(null);
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [index, setIndex] = useState(0);
  const [selected, setSelected] = useState<number | null>(null);
  const [correctIndex, setCorrectIndex] = useState<number | null>(null);
  const [checking, setChecking] = useState(false);

  useEffect(() => {
    let cancelled = false;
    quizApi
      .start(categoryId, size)
      .then((started) => {
        if (cancelled) return;
        setQuiz(started);
        setStatus("running");
      })
      .catch((cause: unknown) => {
        if (cancelled) return;
        setError(
          cause instanceof ApiError ? cause.message : "не вышло начать урок",
        );
        setStatus("failed");
      });
    return () => {
      cancelled = true;
    };
  }, [categoryId, size]);

  const question = quiz?.questions[index] ?? null;
  const total = quiz?.questions.length ?? 0;
  const checked = correctIndex !== null;

  const check = useCallback(async () => {
    if (!quiz || !question || selected === null || checked || checking) return;
    setChecking(true);
    try {
      const result = await quizApi.answer(quiz.id, question.id, selected);
      setQuiz(result.quiz);
      setCorrectIndex(result.correct_index);
      if (!result.correct) {
        learnerStore.loseLife();
      }
    } catch (cause: unknown) {
      setError(
        cause instanceof ApiError ? cause.message : "ответ не сохранился",
      );
    } finally {
      setChecking(false);
    }
  }, [checked, checking, question, quiz, selected]);

  const next = useCallback(() => {
    if (!quiz) return;
    setSelected(null);
    setCorrectIndex(null);
    if (index + 1 >= quiz.questions.length) {
      learnerStore.completeLesson(quiz.category_id, quiz.score.correct);
      setStatus("finished");
      return;
    }
    setIndex(index + 1);
  }, [index, quiz]);

  return {
    status,
    error,
    quiz,
    question,
    index,
    total,
    selected,
    correctIndex,
    checked,
    checking,
    select: setSelected,
    check,
    next,
    /** Доля пройденных вопросов для полосы прогресса. */
    percent: total === 0 ? 0 : Math.round((index / total) * 100),
  };
}
